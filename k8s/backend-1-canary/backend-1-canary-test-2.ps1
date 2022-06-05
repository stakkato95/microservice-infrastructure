
while(1) {
    curl http://localhost/frontend/request -sS -H "X-Api-User-Id: 3" | jq ".nested.nested.service"    
    sleep(0.5);
}