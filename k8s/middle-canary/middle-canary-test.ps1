
while(1) {
    curl http://localhost/frontend/request -sS | jq ".nested.service"    
    sleep(0.5);
}