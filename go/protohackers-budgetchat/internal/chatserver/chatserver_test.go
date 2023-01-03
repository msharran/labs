package chatserver

import (
	"budgetchat/internal/chatserver/mocks"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"
)

//go:generate mockery --name=Conn --srcpkg=net --case=underscore --log-level=info --disable-version-string

func TestChatting(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout))

	s := NewServer(ServerArgs{
		Log:     log,
		Timeout: 5 * time.Millisecond,
	})
	c := s.(*chatServer)
	mockConn := mocks.NewConn(t)

	mockConn.On("SetDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
	mockConn.On("Close").Return(nil).Once()
	readCall := mockConn.On("Read", mock.AnythingOfType("[]uint8")).Return(4, nil).Run(func(args mock.Arguments) {
		buf := args.Get(0).([]byte)
		buf = append(buf, "foo\n"...)
	})
	mockConn.On("Write", []byte("foo\n")).Return(4, nil).Once().NotBefore(readCall)

	c.handleChat(mockConn)
}
