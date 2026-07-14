> 해당 기능은 비공식이며 , 커뮤니티 기여로 작성하는 도구입니다. NHN Dooray 서비스에서 제공하는 기능이 아님을 밝혀 둡니다.

# dooray_mcp

Dooray! 를 Claude 등 MCP 호환 AI 클라이언트에서 사용할 수 있도록 해 주는 **MCP (Model Context Protocol) 서버**입니다.

자연어로 메신저를 보내고, 캘린더 일정을 조회·등록하고, 업무(프로젝트 포스트)를 검색하고, 멤버 정보를 찾을 수 있습니다.

## 주요 기능

| 분류 | 기능 | 관련 도구 |
|------|------|-----------|
| 메신저 | 다른 멤버에게 DM 전송 | `dooray_messenger` |
| 캘린더 | 내 캘린더 목록 조회 | `dooray_calendar_calendars` |
| 캘린더 | 기간별 일정 조회 | `dooray_calendar_events` |
| 캘린더 | 일정 등록 (종일/반복일정 지원) | `dooray_calendar_post_event` |
| 계정 | 이름/userCode 로 멤버 검색 | `dooray_account_members` |
| 계정 | 멤버 상세정보 조회 | `dooray_account_member` |
| 프로젝트 | 참여 중인 프로젝트 조회 | `dooray_project` |
| 프로젝트 | 업무(포스트) 검색 (담당자/상태/기한 필터) | `dooray_posts` |
| 기타 | 현재 시각 조회 | `os` |

반복 일정은 `daily / weekly / monthly / yearly` 주기, interval, 종료일, 요일/일자 지정까지 지원합니다.

## 설치하기

### Homebrew (macOS / Linux)

```bash
brew tap dooray-go/tap
brew install dooray-mcp
```

설치 후 `dooray-mcp` 명령어를 사용할 수 있습니다.

### 직접 다운로드

릴리즈 페이지에서 PC 아키텍처에 해당하는 바이너리를 다운로드 합니다.

* https://github.com/dooray-go/dooray_mcp/releases

### 소스에서 빌드

Go 1.26 이상이 필요합니다.

```bash
git clone https://github.com/dooray-go/dooray_mcp.git
cd dooray_mcp

# 현재 플랫폼용 빌드
go build -o dist/dooray-mcp .

# 또는 모든 플랫폼용 크로스 컴파일
make build-all
```

`make build-all` 실행 시 `dist/` 디렉터리에 다음 바이너리가 생성됩니다.

* `dooray.darwin.amd64`, `dooray.darwin.arm64`
* `dooray.linux.amd64`
* `dooray.windows.amd64.exe`

## Dooray 개인 인증 토큰 발급

1. Dooray! 웹에서 **개인설정 > API > 개인 인증 토큰** 메뉴로 이동합니다.
2. 새 토큰을 생성하고 복사해 둡니다. (토큰은 한 번만 노출됩니다.)
3. 아래 설정에서 `{개인토큰}` 자리에 이 값을 입력합니다.

## 설정하기

### Claude Desktop 에서 사용하기

1. MCP 를 지원하는 Claude.ai 데스크탑 애플리케이션을 설치합니다. → https://claude.ai/download
2. Claude.ai 데스크탑 애플리케이션에서 **설정 > 개발자 > 설정 편집** 을 선택합니다.

   ![img.png](img.png)

3. `claude_desktop_config.json` 을 다음과 같이 편집합니다.

```json
{
  "mcpServers": {
    "dooray": {
      "command": "dooray-mcp",
      "args": [
        "--token",
        "{개인토큰}"
      ]
    }
  },
  "globalShortcut": ""
}
```

> Homebrew 로 설치한 경우 `dooray-mcp` 를 그대로 사용합니다. 직접 다운로드한 경우 바이너리의 전체 경로를 입력하세요.
> 예: `"/Users/me/bin/dooray.darwin.arm64"`

4. Claude Desktop 을 재시작하면 Dooray 도구가 인식됩니다.

### Claude Code (CLI) 에서 사용하기

Claude Code 에 MCP 서버로 등록하면 터미널에서 바로 Dooray 기능을 사용할 수 있습니다.

1. 프로젝트 단위 등록

```bash
claude mcp add dooray -- dooray-mcp --token {개인토큰}
```

2. 전역으로 등록

```bash
claude mcp add --scope user dooray -- dooray-mcp --token {개인토큰}
```

3. 등록 확인

```bash
claude mcp list
```

4. 사용 예

```bash
claude "오늘 내 캘린더 일정을 알려줘"
```

## 사용 예시

### 메신저

```
오늘 내 일정을 중요한 순으로 정렬해서 김XX 에게 메신저로 보내 줘.
```

```
정만티에게 "회의 시작합니다" 라고 DM 보내줘.
```

### 캘린더 조회

```
내일 일정 중에 중요한 일정은 뭐야?
```

```
이번 주 금요일 오후에 비어있는 시간 알려줘.
```

### 일정 등록

```
내일 아침 9시에 자유수영 할 예정이야. 두시간 일정 등록해줘.
```

```
다음주부터 매주 월,수,금 오전 10시 스크럼 30분 일정 등록해줘. 4주간 반복.
```

```
매달 1일 오전 10시에 월간 리뷰 일정 등록, 올해 12월까지.
```

![img_1.png](img_1.png)

### 업무 / 프로젝트

```
나는 정만티야 Dooray-잘쓰자 프로젝트에서 가장 급한일을 알려줘
```

```
Dooray-잘쓰자 프로젝트에서 내게 할당된 업무 중 이번 주 마감인 것만 알려줘.
```

```
지난 30일간 생성된 내 업무를 상태별로 정리해 줘.
```

## 도구 레퍼런스

### `dooray_messenger`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `send` |
| to | O | 수신자의 organizationMemberId |
| message | O | 보낼 메시지 본문 |

### `dooray_calendar_calendars`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_calendars` |

### `dooray_calendar_events`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_events` |
| calendars | X | 조회할 캘린더 ID 목록 (쉼표 구분) |
| timeMin | O | 시작 시각 (ISO 8601, 예: `2025-04-11T00:00:00+09:00`) |
| timeMax | O | 종료 시각 (ISO 8601) |

### `dooray_calendar_post_event`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `create_event` |
| calendarId | X | 등록 대상 캘린더 ID |
| subject | O | 일정 제목 |
| content | O | 일정 본문 (text/html) |
| startedAt | O | 시작 시각 (ISO 8601) |
| endedAt | O | 종료 시각 (ISO 8601) |
| wholeDayFlag | X | 종일 일정 여부. `true` 인 경우 날짜만 지정 (예: `2025-04-11+09:00`) |
| recurrenceFrequency | X | `daily` / `weekly` / `monthly` / `yearly` |
| recurrenceInterval | X | 반복 간격 (기본 1) |
| recurrenceUntil | X | 반복 종료일 (ISO 8601) |
| recurrenceByday | X | 반복 요일. 예: `MO,WE,FR`, 월간은 `1MO`, `-1FR` 등 |
| recurrenceBymonth | X | 반복 월(1-12), 연간 반복용 |
| recurrenceBymonthday | X | 반복 일(1-31), 월간/연간 반복용 |
| recurrenceTimezoneName | X | 타임존 (기본 `Asia/Seoul`) |

### `dooray_account_members` / `dooray_account_member`

| 도구 | operation | 설명 |
|------|-----------|------|
| `dooray_account_members` | `find_member_id` | 이름 또는 userCode 로 멤버 검색 |
| `dooray_account_member` | `find_member_details` | memberId 로 상세정보(닉네임, 이름 등) 조회 |

### `dooray_project`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_projects` |
| type | O | `public` / `private` |
| state | O | `active` / `archived` |
| scope | O | `private` / `public` |

### `dooray_posts`

업무(포스트) 검색 도구. `projectId` 만 필수이며, 나머지는 필터로 사용됩니다.

| 파라미터 | 설명 |
|----------|------|
| operation | `find_posts` (필수) |
| projectId | 프로젝트 ID (필수, 쉼표로 여러 개 지정 가능) |
| page / size | 페이지(기본 0), 페이지 크기(기본 20, 최대 100) |
| fromEmailAddress | 보낸 사람 이메일로 필터 |
| fromMemberIds | 작성자 memberId (쉼표 구분) |
| toMemberIds | 담당자 memberId |
| toMemberSize | 담당자 수 (0: 미지정, 1: 단일 담당자) |
| ccMemberIds | 참조자 memberId |
| tagIds | 태그 ID |
| parentPostId | 상위 업무 ID (하위 업무 조회) |
| postNumber | 업무 번호 |
| postWorkflowClasses | `backlog` / `registered` / `working` / `closed` |
| postWorkflowIds | 워크플로 ID |
| milestoneIds | 마일스톤 ID |
| subjects | 제목 키워드 |
| createdAt / updatedAt / dueAt | 날짜 필터. `today`, `thisweek`, `prev-30d`, `next-7d`, 또는 ISO8601 구간 `~` 형식 |
| order | 정렬: `postDueAt`, `postUpdatedAt`, `createdAt` (내림차순은 `-` 접두사) |

### `os`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `get_date_time` — 현재 로컬 시각 반환 |

## 개발

### 의존성

* [github.com/dooray-go/dooray](https://github.com/dooray-go/dooray) — Dooray OpenAPI Go 클라이언트
* [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) — MCP 서버 SDK

### 디렉터리 구조

```
.
├── main.go             # MCP 서버 초기화 및 도구 등록
├── os.go               # 시간 관련 도구
├── account.go          # 멤버 조회 도구
├── calendar.go         # 캘린더/일정 도구 (반복 일정 포함)
├── messenger.go        # 메신저 DM 도구
├── project.go          # 프로젝트/업무 도구
├── Makefile            # 크로스 컴파일 빌드 스크립트
└── *_test.go           # 각 도구별 단위 테스트
```

### 테스트

```bash
go test ./...
```

### 빌드

```bash
make build-all   # 모든 타겟(darwin/linux/windows) 빌드
make clean       # dist/ 제거
```

## 문제 해결

* **`token must be set!!` 로그 후 종료**: `--token` 인자가 지정되지 않았습니다. 설정 파일이나 CLI 인자를 확인하세요.
* **Claude Desktop 에서 도구가 보이지 않음**: 설정 편집 후 Claude Desktop 을 완전히 종료했다가 다시 실행해야 합니다.
* **일정 등록 시 시간 파싱 오류**: `startedAt`, `endedAt` 은 반드시 ISO 8601 형식이어야 합니다 (예: `2025-04-11T09:00:00+09:00`). 종일 일정은 `2025-04-11+09:00` 형태로 지정합니다.

## 라이선스 / 기여

이슈와 PR 은 언제나 환영합니다. → https://github.com/dooray-go/dooray_mcp
