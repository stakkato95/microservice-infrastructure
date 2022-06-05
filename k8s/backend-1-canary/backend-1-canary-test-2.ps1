
while(1) {
    curl http://localhost/request -sS -H "X-Api-User-Id: 3" | jq ".nested.nested.service"    
    sleep(0.5);
}