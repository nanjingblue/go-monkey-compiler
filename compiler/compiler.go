package compiler

import (
	"github.com/nanjingblue/go-monkey/ast"
	"github.com/nanjingblue/go-monkey/code"
	"github.com/nanjingblue/go-monkey/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object // 常量池
}

/*
Bytecode 是需要传输给虚拟机并在编译器测试中做断言的内容
*/
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

/*
addConstant 将 obj 添加到 constants 切片的末尾，通过返回其在 constants 切片中的索引来引来为其提供标识符
此标识符现在将用作 OpConstant 指令的操作数，该指令驱使虚拟机从常量池加载此常量至栈
*/
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

/*
emit "发出"（“生成” / "输出"）指令
生成指令并将其添加到最终结果
*/
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}
