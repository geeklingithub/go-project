package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/**
掌握context.WithXXXX 的几个 API 使用方式
*/

/**
掌握 Context 上的 Err, Done, 和 Value 的用法
*/

// Background returns a non-nil, empty Context. It is never canceled, has no
// values, and has no deadline. It is typically used by the main function,
// initialization, and tests, and as the top-level Context for incoming
// requests.
func TestContextBackground(t *testing.T) {
	ctx := context.Background()
	fmt.Println(ctx)

	fmt.Println(ctx.Err())
	fmt.Println(ctx.Done())
	fmt.Println(ctx.Deadline())
	fmt.Println(ctx.Value("key"))

}

// TODO returns a non-nil, empty Context. Code should use context.TODO when
// it's unclear which Context to use or it is not yet available (because the
// surrounding function has not yet been extended to accept a Context
// parameter).
func TestContextTODO(t *testing.T) {
	ctx := context.TODO()
	fmt.Println(ctx)
}

// WithDeadline returns a copy of the parent context with the deadline adjusted
// to be no later than d. If the parent's deadline is already earlier than d,
// WithDeadline(parent, d) is semantically equivalent to parent. The returned
// context's Done channel is closed when the deadline expires, when the returned
// cancel function is called, or when the parent context's Done channel is
// closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func TestContextWithDeadline(t *testing.T) {
	ctx := context.Background()
	dlCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second))

	// Even though ctx will be expired, it is good practice to call its
	// cancellation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("overslept")
	case <-dlCtx.Done():
		// 区分主动取消或超时
		fmt.Println("vCtx.Err():", dlCtx.Err())
		//channel通信协调，配合select-case使用
		fmt.Println("vCtx.Done():", dlCtx.Done())
		//获取数据
		fmt.Println("vCtx.Value():", dlCtx.Value("key"))
	}
}

func TestContextWithTimeout(t *testing.T) {
	ctx := context.Background()
	dlCtx, cancel := context.WithTimeout(ctx, time.Second)

	defer cancel()

	select {
	case <-time.After(1 * time.Millisecond):
		fmt.Println("overslept")
	case <-dlCtx.Done():

		fmt.Println(dlCtx.Err())
	}
}

func TestContextCancel(t *testing.T) {
	ctx := context.Background()
	dlCtx, cancel := context.WithCancel(ctx)
	cancel()
	fmt.Println(dlCtx.Value("KEY"))
	fmt.Println(dlCtx.Err())
}

func TestContextWithValue(t *testing.T) {
	ctx := context.Background()
	vCtx := context.WithValue(ctx, "key", "value")
	value := vCtx.Value("key")
	fmt.Println(value)
}

func TestContextErr(t *testing.T) {
	ctx := context.Background()
	vCtx := context.WithValue(ctx, "key", "value")
	fmt.Println("vCtx.Err():", vCtx.Err())
	fmt.Println("vCtx.Done():", vCtx.Done())
	fmt.Println("vCtx.Value():", vCtx.Value("key"))
}

func TestWithCancels(t *testing.T) {
	c1, cancel := context.WithCancel(context.Background())

	if got, want := fmt.Sprint(c1), "context.Background.WithCancel"; got != want {
		t.Errorf("c1.String() = %q want %q", got, want)
	}

	o := context.Background()
	c2, _ := context.WithCancel(o)
	contexts := []context.Context{c1, o, c2}

	for i, c := range contexts {
		if d := c.Done(); d == nil {
			t.Errorf("c[%d].Done() == %v want non-nil", i, d)
		}
		if e := c.Err(); e != nil {
			t.Errorf("c[%d].Err() == %v want nil", i, e)
		}

		select {
		case x := <-c.Done():
			t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
		default:
		}
	}

	cancel() // Should propagate synchronously.
	for i, c := range contexts {
		select {
		case <-c.Done():
			fmt.Println(i, c)
		default:
			//t.Errorf("<-c[%d].Done() blocked, but shouldn't have", i)
		}
		//if e := c.Err(); e != context.Canceled {
		//	t.Errorf("c[%d].Err() == %v want %v", i, e, context.Canceled)
		//}
	}
}

type otherContext struct {
	context.Context
}
