import json
import os

collection_path = r'd:\pg management system\postman_collection.json'

with open(collection_path, 'r') as f:
    data = json.load(f)

def update_item(item):
    if 'request' in item:
        url = item['request']['url']
        path = url.get('path', [])
        
        # Don't prefix health, auth, graphql, or debug routes
        excluded_top_paths = ['health', 'auth', 'graphql', 'debug']
        
        should_prefix = False
        if path:
            if path[0] not in excluded_top_paths and path[0] != 'api':
                should_prefix = True
        
        if should_prefix:
            url['path'] = ['api'] + path
            # Update raw URL string
            url['raw'] = url['raw'].replace('localhost:8080/', 'localhost:8080/api/')
            
            # Add Authorization header if it doesn't exist
            if 'header' not in item['request']:
                item['request']['header'] = []
            
            has_auth = any(h.get('key') == 'Authorization' for h in item['request']['header'])
            if not has_auth:
                item['request']['header'].append({
                    "key": "Authorization",
                    "value": "Bearer {{token}}",
                    "type": "text"
                })
                
    if 'item' in item:
        for sub_item in item['item']:
            update_item(sub_item)

for item in data['item']:
    update_item(item)

# Add collection variable for token
if 'variable' not in data:
    data['variable'] = []

# Update or add token variable
token_found = False
for var in data['variable']:
    if var.get('key') == 'token':
        token_found = True
        break

if not token_found:
    data['variable'].append({
        "key": "token",
        "value": "",
        "type": "string"
    })

with open(collection_path, 'w') as f:
    json.dump(data, f, indent=4)

print("Postman collection updated successfully.")
