import json

with open("postman/postman_collection.json", "r") as f:
    data = json.load(f)

new_folder = {
    "name": "ExamRoom Management",
    "item": [
        {
            "name": "Get List Exam Rooms",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {{admin_access_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/v1/admin/exam-rooms?page=1&limit=10&keyword=",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "admin", "exam-rooms"],
                    "query": [
                        {"key": "page", "value": "1"},
                        {"key": "limit", "value": "10"},
                        {"key": "keyword", "value": ""}
                    ]
                }
            }
        },
        {
            "name": "Get Exam Room Detail",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {{admin_access_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/v1/admin/exam-rooms/1",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "admin", "exam-rooms", "1"]
                }
            }
        },
        {
            "name": "Create Exam Room",
            "request": {
                "method": "POST",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {{admin_access_token}}"
                    },
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\n    \"room_name\": \"Phòng Thi 101\",\n    \"location\": \"Tầng 1 - Tòa A\",\n    \"capacity\": 30\n}"
                },
                "url": {
                    "raw": "{{base_url}}/api/v1/admin/exam-rooms",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "admin", "exam-rooms"]
                }
            }
        },
        {
            "name": "Update Exam Room",
            "request": {
                "method": "PUT",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {{admin_access_token}}"
                    },
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\n    \"room_name\": \"Phòng Thi 101 Cập Nhật\",\n    \"location\": \"Tầng 1 - Tòa A\",\n    \"capacity\": 40\n}"
                },
                "url": {
                    "raw": "{{base_url}}/api/v1/admin/exam-rooms/1",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "admin", "exam-rooms", "1"]
                }
            }
        },
        {
            "name": "Delete Exam Room",
            "request": {
                "method": "DELETE",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {{admin_access_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/v1/admin/exam-rooms/1",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "admin", "exam-rooms", "1"]
                }
            }
        }
    ]
}

exists = False
for item in data['item']:
    if item['name'] == 'ExamRoom Management':
        item['item'] = new_folder['item']
        exists = True
        break

if not exists:
    data['item'].append(new_folder)

with open("postman/postman_collection.json", "w") as f:
    json.dump(data, f, indent=4, ensure_ascii=False)
