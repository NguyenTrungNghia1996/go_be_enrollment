import json

with open("postman/postman_collection.json", "r") as f:
    data = json.load(f)

new_folder = {
    "name": "AdmissionPeriodSubject Management",
    "item": [
        {
            "name": "Get Admission Period Subjects",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {{admin_access_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/v1/admin/admission-periods/1/subjects",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "admin", "admission-periods", "1", "subjects"]
                }
            }
        },
        {
            "name": "Update Admission Period Subjects",
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
                    "raw": "{\n    \"subjects\": [\n        {\n            \"subject_id\": 1,\n            \"weight\": 2.0,\n            \"is_required\": true\n        },\n        {\n            \"subject_id\": 2,\n            \"weight\": 1.0,\n            \"is_required\": false\n        }\n    ]\n}"
                },
                "url": {
                    "raw": "{{base_url}}/api/v1/admin/admission-periods/1/subjects",
                    "host": ["{{base_url}}"],
                    "path": ["api", "v1", "admin", "admission-periods", "1", "subjects"]
                }
            }
        }
    ]
}

exists = False
for item in data['item']:
    if item['name'] == 'AdmissionPeriodSubject Management':
        item['item'] = new_folder['item']
        exists = True
        break

if not exists:
    data['item'].append(new_folder)

with open("postman/postman_collection.json", "w") as f:
    json.dump(data, f, indent=4, ensure_ascii=False)
