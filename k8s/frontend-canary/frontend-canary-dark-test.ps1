
while(1) {
    curl http://localhost/request -sS -H "X-Api-User-Id: 4" -H "X-New-Frontend: 100500" | jq ".service"    
    sleep(0.5);
}