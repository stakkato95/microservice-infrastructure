
while(1) {
    curl http://localhost/frontend/request -sS -H "X-Api-User-Id: 4" | jq ".nested.nested.service"    
    sleep(0.5);
}