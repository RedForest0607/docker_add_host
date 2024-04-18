package utils

import "regexp"

func IsValidIPAddress(ip string) bool {
	// IP 주소 패턴 정의
	ipPattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`

	// 정규 표현식으로 주어진 문자열과 패턴을 비교하여 일치 여부 확인
	match, _ := regexp.MatchString(ipPattern, ip)

	return match
}
