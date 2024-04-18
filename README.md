## Docker Container의 /etc/hosts에 새로운 정보 기입과, docker-compose 파일에 extra_hosts에 정보를 기입하는 유틸리티 입니다.  

`./bin/dockerAddHost`

다음과 같은 추가 기능안이 있습니다 
- [ ] 프롬프트를 통한 방식이 아닌, 명령어 한줄로 실행 가능한 형태의 코드
- [ ] 이미 존재하는 host 여부 확인
- [ ] host 정보 삭제 확인
- [ ] host name이 존재한다면 이어 붙이기
- [ ] 컨테이너 적용, yaml 적용 중 하나가 오류가 발생할 시 트랜젝션 처리 (bak 파일 형식으로 처리해야할 것으로 보임)