# Ghost Images

Ghost 블로그 플랫폼에 이미지를 업로드하고 관리하기 위한 API 서버입니다. 이 프로젝트는 Ghost Admin API를 활용하여 이미지 업로드 및 관리 기능을 제공합니다.
기존 Ghost MCP에서 지원하지 않는 기능을 추가합니다. 
## 목적 
node python 기반의 mcp 서버들이 많은데 현재 저는 60가지 이상의 툴을 사용하고 있습니다. 
이러한 서버들이 많은 메모리를 차지 하게 되었으며 go native 기반의 mcp 서버를 사용하면 메모리 사용량을 줄일 수 있을 것으로 예상합니다.
## 기능

- Ghost 블로그에 이미지 업로드
- local path 이미지 업로드
- Base64 file save after upload
- JWT 인증
- MCP 서버 기능

## 요구사항

- Go 1.24 이상
- Ghost 블로그 인스턴스
- Ghost Admin API 키

## 설치 방법

1. 저장소 클론

```bash
git clone https://github.com/gudcks0305/ghost-mcp-images.git
cd ghost-mcp-images
```

2. 의존성 설치

```bash
go mod download
```

3. 환경 변수 설정

`.env.example` 파일을 복사하여 `.env` 파일을 생성하고 필요한 설정을 입력합니다.

```bash
cp .env.example .env
```

`.env` 파일을 편집하여 Ghost API URL과 API 키를 설정합니다:

```
GHOST_API_URL=https://your-ghost-instance.com
GHOST_STAFF_API_KEY=your-staff-api-key
```

Ghost Admin API 키는 Ghost 관리자 패널의 설정 > 통합 메뉴에서 찾을 수 있습니다.

## 빌드 방법

```bash
go build -o bin/ghost-images ./cmd/main.go
```

## 실행 방법

### 직접 실행

```bash
go run cmd/main.go
```

### 빌드 후 실행

```bash
./bin/ghost-images
```

### Air를 사용한 개발 모드 실행 (라이브 리로딩)

```bash
air
```

## API 사용 예제

### 이미지 업로드

```
내부적으로 curl 명령어를 사용하여 이미지를 업로드합니다. http 415 error 발생으로 아래 curl 명령어를 사용하게 되었습니다.
curl -X POST http://localhost:8080/upload -H "Content-Type: multipart/form-data" -F "file=@/path/to/image.jpg"
```
실제 요청은 JSON-RPC 형식으로 요청합니다.



## 개발 팁

1. `.air.toml` 파일을 수정하여 라이브 리로딩 설정을 변경할 수 있습니다.
2. 환경 변수는 `.env` 파일 대신 시스템 환경 변수로 설정할 수도 있습니다.

## 라이센스

[MIT](LICENSE)

## 연락처

문의사항이 있으시면 이슈를 생성하거나 이메일로 연락해주세요.

## Claude Desktop Config
### mac 환경에서는 `~/Library/Application Support/Claude/claude_desktop_config.json` 경로에 위 설정을 추가합니다.

```
{
    "ghost-images": {
      "command": "zsh",
      "args": [
        "-c",
        "your-builded-path-example/bin/ghost-images/main"
      ],
      "env": {
        "GHOST_API_URL": "https://your-ghost-instance.com",
        "GHOST_STAFF_API_KEY": "your-staff-api-key"
      }
    }
}
```
### windows 환경에서는 `C:\Users\<username>\AppData\Roaming\anthropic\claude_desktop_config.json` 경로에 위 설정을 추가합니다.  

```
{
    "ghost-images": {
      "command": "cmd",
      "args": [
        "/c",
        "your-builded-path-example\\bin\\ghost-images\\main.exe"
      ],
      "env": {
        "GHOST_API_URL": "https://your-ghost-instance.com",
        "GHOST_STAFF_API_KEY": "your-staff-api-key"
      }
    }
}
```

## 참고
- [Go Mcp](https://github.com/mark3labs/mcp-go)
- [Ghost Admin API](https://ghost.org/docs/api/admin/)
- [JSON-RPC 2.0](https://www.jsonrpc.org/specification)
- [MFYDev/ghost-mcp](https://github.com/MFYDev/ghost-mcp)

### Claude Desktop 프롬프트 사용 예시
![image](image.png)