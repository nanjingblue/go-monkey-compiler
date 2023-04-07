package vm

import (
	"fmt"
	"github.com/nanjingblue/go-monkey/code"
	"github.com/nanjingblue/go-monkey/compiler"
	"github.com/nanjingblue/go-monkey/object"
)

const StackSize = 2048

/*
VM 虚拟机是一个有 4 个字段的结构体
包含编译器 compiler 生成的常量 constants 和指令 instructions
一个栈 stack stack 预分配了 StackSize 数量的元素
sp 是栈指针 通过递增或递减来增大或缩小栈 而不是通过修改栈切片本身
*/
type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // 始终指向栈中的下一个空间槽，栈顶的值为 stack[sp - 1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// StackTop 返回栈顶元素
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

/*
LastPoppedStackElem 获取最后一个弹出栈的元素
由于只是通过递减 vm.sp 使元素弹栈 因此可以找到之前栈顶元素的地方
修改虚拟机测试，确保 OpPop 的执行时正确的
*/
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operation: %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}