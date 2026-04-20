package service

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"go_be_enrollment/internal/common/storage"
	"go_be_enrollment/internal/modules/applicationdocument/dto"
	"go_be_enrollment/internal/modules/applicationdocument/entity"
	"go_be_enrollment/internal/modules/applicationdocument/repository"
)

type ApplicationDocumentService interface {
	GetAdminList(appID uint) ([]dto.ApplicationDocumentRes, error)
	GetUserList(appID, userID uint) ([]dto.ApplicationDocumentRes, error)
	UploadDocument(appID, userID uint, docType string, fileHeader *multipart.FileHeader) (*dto.ApplicationDocumentRes, error)
	DeleteDocument(id, userID uint) error
}

type applicationDocumentService struct {
	repo    repository.ApplicationDocumentRepository
	storage storage.StorageService
}

func NewApplicationDocumentService(repo repository.ApplicationDocumentRepository, storageService storage.StorageService) ApplicationDocumentService {
	return &applicationDocumentService{
		repo:    repo,
		storage: storageService,
	}
}

func (s *applicationDocumentService) mapToDto(doc *entity.ApplicationDocument) dto.ApplicationDocumentRes {
	var publicURL string
	if s.storage != nil {
		publicURL = s.storage.GetPublicURL(doc.FilePath)
	}

	return dto.ApplicationDocumentRes{
		ID:            doc.ID,
		ApplicationID: doc.ApplicationID,
		DocumentType:  doc.DocumentType,
		FilePath:      doc.FilePath,
		PublicURL:     publicURL,
		CreatedAt:     doc.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *applicationDocumentService) ValidateFile(fileHeader *multipart.FileHeader) error {
	// Giới hạn 5MB
	const MAX_SIZE = 5 * 1024 * 1024
	if fileHeader.Size > MAX_SIZE {
		return errors.New("kích thước file phải nhỏ hơn 5MB")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{
		".pdf":  true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}

	if !allowedExts[ext] {
		return errors.New("hệ thống chỉ hỗ trợ file PDF, PNG, JPG, JPEG")
	}

	return nil
}

func (s *applicationDocumentService) GetAdminList(appID uint) ([]dto.ApplicationDocumentRes, error) {
	list, err := s.repo.GetByApplicationID(appID)
	if err != nil {
		return nil, errors.New("lỗi lấy danh sách tài liệu")
	}

	res := make([]dto.ApplicationDocumentRes, 0)
	for _, doc := range list {
		res = append(res, s.mapToDto(&doc))
	}
	return res, nil
}

func (s *applicationDocumentService) GetUserList(appID, userID uint) ([]dto.ApplicationDocumentRes, error) {
	isOwner, _, err := s.repo.VerifyApplicationAccess(appID, userID)
	if err != nil {
		return nil, errors.New("lỗi hệ thống khi xác thực hồ sơ")
	}
	if !isOwner {
		return nil, errors.New("không tìm thấy hồ sơ của bạn")
	}

	list, err := s.repo.GetByApplicationID(appID)
	if err != nil {
		return nil, errors.New("lỗi lấy danh sách tài liệu")
	}

	res := make([]dto.ApplicationDocumentRes, 0)
	for _, doc := range list {
		res = append(res, s.mapToDto(&doc))
	}
	return res, nil
}

func (s *applicationDocumentService) UploadDocument(appID, userID uint, docType string, fileHeader *multipart.FileHeader) (*dto.ApplicationDocumentRes, error) {
	isOwner, status, err := s.repo.VerifyApplicationAccess(appID, userID)
	if err != nil {
		return nil, errors.New("lỗi hệ thống khi xác thực hồ sơ")
	}
	if !isOwner {
		return nil, errors.New("không tìm thấy hồ sơ của bạn")
	}
	if status != "Draft" {
		return nil, errors.New("chỉ có thể upload tài liệu khi hồ sơ ở trạng thái nháp")
	}

	if err := s.ValidateFile(fileHeader); err != nil {
		return nil, err
	}

	objectKey := s.storage.BuildObjectKey(appID, fileHeader.Filename)

	if err := s.storage.UploadFile(fileHeader, objectKey); err != nil {
		return nil, errors.New("không thể tải file lên hệ thống")
	}

	doc := &entity.ApplicationDocument{
		ApplicationID: appID,
		DocumentType:  docType,
		FilePath:      objectKey,
	}

	if err := s.repo.Create(doc); err != nil {
		// Rollback upload since DB failed
		_ = s.storage.DeleteFile(objectKey)
		return nil, errors.New("lỗi khi lưu thông tin tài liệu")
	}

	res := s.mapToDto(doc)
	return &res, nil
}

func (s *applicationDocumentService) DeleteDocument(id, userID uint) error {
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy tài liệu")
	}

	isOwner, status, err := s.repo.VerifyApplicationAccess(doc.ApplicationID, userID)
	if err != nil {
		return errors.New("lỗi hệ thống khi xác thực hồ sơ")
	}
	if !isOwner {
		return errors.New("bạn không có quyền xóa tài liệu này")
	}
	if status != "Draft" {
		return errors.New("chỉ có thể xóa tài liệu khi hồ sơ ở trạng thái nháp")
	}

	// Xóa DB trước, tránh DB lỗi mà file đã bị xóa
	if err := s.repo.Delete(id); err != nil {
		return errors.New("lỗi khi xóa tài liệu khỏi hệ thống")
	}

	// Ignore physical delete fail
	_ = s.storage.DeleteFile(doc.FilePath)

	return nil
}
