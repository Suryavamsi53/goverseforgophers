# Mocking and Interfaces

Unit Tests must be fast (executing in <1ms). They must be deterministic (not failing because the wifi dropped). 

If your `UserService` directly imports the `database/sql` package to query Postgres, you cannot write a Unit Test. You would have to boot up a real Postgres database, seed it with test data, and ensure port 5432 is open. That is an Integration Test.

To write true Unit Tests, you must use **Mocking**.

## 1. The Interface Boundary

As we learned in the Clean Architecture module ("Accept Interfaces, Return Structs"), Mocking is only possible if your Service accepts an Interface.

```go
// 1. The Consumer defines the Interface
type UserRepository interface {
    GetUser(id int) (*User, error)
}

// 2. The Service requires the Interface via Constructor Injection
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

## 2. Writing a Manual Mock

In many languages, you have to install massive frameworks to generate Mocks. 
In Go, because interfaces are satisfied implicitly, you just write a dummy struct in your test file!

```go
// In user_service_test.go

// 1. Create a dummy struct
type MockUserRepo struct {
    // We add fields to allow the test to control the mock's behavior!
    mockData *User
    mockErr  error
}

// 2. Implement the Interface
func (m *MockUserRepo) GetUser(id int) (*User, error) {
    return m.mockData, m.mockErr
}
```

## 3. Injecting the Mock

Now, testing the `UserService` is completely deterministic. We don't need Postgres. We just inject the Mock, set the desired outcome, and assert the behavior!

```go
func TestUserService_Success(t *testing.T) {
    // Inject the mock with fake data
    mockRepo := &MockUserRepo{
        mockData: &User{ID: 1, Name: "Test User"},
    }
    
    svc := NewUserService(mockRepo)
    
    user, err := svc.GetProfile(1)
    if err != nil { t.Errorf("Expected success, got %v", err) }
    if user.Name != "Test User" { t.Errorf("Got %s", user.Name) }
}

func TestUserService_DatabaseCrash(t *testing.T) {
    // Inject the mock and force it to simulate a DB crash!
    mockRepo := &MockUserRepo{
        mockErr: errors.New("connection refused"),
    }
    
    svc := NewUserService(mockRepo)
    
    _, err := svc.GetProfile(1)
    if err == nil { t.Error("Expected error due to DB crash, got nil") }
}
```

## 4. GoMock (Enterprise Mock Generation)

Writing manual mocks is great for a 2-method interface. If your `UserRepository` has 45 methods, writing manual mocks is agonizing.

Enterprise Go teams use **GoMock** (built by Google, now maintained by Uber).

GoMock is a CLI tool (`mockgen`) that reads your interfaces and automatically generates thousands of lines of Mock code.

```bash
mockgen -source=user_service.go -destination=mocks/mock_user_repo.go
```

In your test, GoMock allows you to set strict expectations on how many times a method should be called!

```go
func TestWithGoMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Instantiate the auto-generated mock
    mockRepo := mocks.NewMockUserRepository(ctrl)

    // Set an Expectation: "I expect GetUser(1) to be called EXACTLY once. 
    // Return a fake user when it happens."
    mockRepo.EXPECT().
        GetUser(1).
        Return(&User{Name: "GoMock"}, nil).
        Times(1)

    svc := NewUserService(mockRepo)
    svc.GetProfile(1)
}
```
If the `UserService` accidentally calls `GetUser(1)` twice, or forgets to call it at all, the test will automatically fail!
