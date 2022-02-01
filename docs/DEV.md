### Go Module Command
```
go mod init [module-name]
명령어의 이름에서도 느껴지듯이 모듈을 처음 사용할때 사용한다. module-name은 보통 github.com/jay/hello 포멧을 취한다.

go get [module-name]
모듈을 다운로드하는 명령어.

go mod tidy
소스 코드를 확인해서 import되지 않는 모듈들을 자동으로 go.mod 파일에서 삭제하고 import되었지만 실제 모듈이 다운안된 경우는 go.mod파일에 추가해준다.

go mod vendor
Module을 이용하면 module 들을 project 밑에 저장하지 않고, GOPATH에 저장하게 된다. 그러나 자신이 이용하던 모듈들을 repo에 넣고 싶을 경우가 있다. 자동으로 변경될수 있는 모듈들을 고정시키고 싶을때 말이다. 물론 버젼을 강제 지정할수도 있지만, 그 패지키 자체를 가지고 있는것도 쉬운 방법. 따라서 이 명령어를 실행시키면 사용하는 모듈들을 자신의 repo 아래 vendor폴더에 복사를 하게 된다.
```

### 참고
```
gin
https://github.com/gin-gonic/gin#parameters-in-path

gorm
https://gorm.io/ko_KR/docs/advanced_query.html#SubQuery

go excelize
https://github.com/qax-os/excelize

```