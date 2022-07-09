package nahi

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"

	lexer "github.com/DonnieTD/NAH/Lexer"
	utils "github.com/DonnieTD/NAH/Utils"
)

type NAH struct {
	LEXER *lexer.Lexer
}

func GenerateAssemblyForDump(datawriter *bufio.Writer) {
	datawriter.WriteString("dump:\n")
	datawriter.WriteString("    mov     r9, -3689348814741910323\n")
	datawriter.WriteString("    sub     rsp, 40\n")
	datawriter.WriteString("    mov     BYTE [rsp+31], 10\n")
	datawriter.WriteString("    lea     rcx, [rsp+30]\n")
	datawriter.WriteString(".L2:\n")
	datawriter.WriteString("    mov     rax, rdi\n")
	datawriter.WriteString("    lea     r8, [rsp+32]\n")
	datawriter.WriteString("    mul     r9\n")
	datawriter.WriteString("    mov     rax, rdi\n")
	datawriter.WriteString("    sub     r8, rcx\n")
	datawriter.WriteString("    shr     rdx, 3\n")
	datawriter.WriteString("    lea     rsi, [rdx+rdx*4]\n")
	datawriter.WriteString("    add     rsi, rsi\n")
	datawriter.WriteString("    sub     rax, rsi\n")
	datawriter.WriteString("    add     eax, 48\n")
	datawriter.WriteString("    mov     BYTE [rcx], al\n")
	datawriter.WriteString("    mov     rax, rdi\n")
	datawriter.WriteString("    mov     rdi, rdx\n")
	datawriter.WriteString("    mov     rdx, rcx\n")
	datawriter.WriteString("    sub     rcx, 1\n")
	datawriter.WriteString("    cmp     rax, 9\n")
	datawriter.WriteString("    ja      .L2\n")
	datawriter.WriteString("    lea     rax, [rsp+32]\n")
	datawriter.WriteString("    mov     edi, 1\n")
	datawriter.WriteString("    sub     rdx, rax\n")
	datawriter.WriteString("    xor     eax, eax\n")
	datawriter.WriteString("    lea     rsi, [rsp+32+rdx]\n")
	datawriter.WriteString("    mov     rdx, r8\n")
	datawriter.WriteString("    mov     rax, 1\n")
	datawriter.WriteString("    syscall\n")
	datawriter.WriteString("    add     rsp, 40\n")
	datawriter.WriteString("    ret\n")
}

func (n *NAH) Compile() {
	utils.CountTokensCheck(lexer.COUNT_TOKENS, 12, "./NAHI/NAHI.go", "Compile")

	if _, err := os.Stat("./" + "output.asm"); err == nil {
		e := os.Remove("output.asm")
		if e != nil {
			log.Fatal(e)
		}
	}

	file, err := os.OpenFile("output.asm", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)
	datawriter.WriteString("segment .text" + "\n")
	GenerateAssemblyForDump(datawriter)
	datawriter.WriteString("global _start" + "\n")
	datawriter.WriteString("_start:" + "\n")

	for token_index, token := range n.LEXER.Tokens {
		datawriter.WriteString(fmt.Sprintf("addr_%d: \n", token_index))
		switch token.TokenType {
		case lexer.TOKEN_PUSH:
			datawriter.WriteString(fmt.Sprintf("    ;;-- push %d --", token.Parameter) + "\n")
			datawriter.WriteString(fmt.Sprintf("    push %d", token.Parameter) + "\n")
		case lexer.TOKEN_PLUS:
			datawriter.WriteString("    ;;-- plus %d -- \n")
			datawriter.WriteString("    pop rax \n")
			datawriter.WriteString("    pop rbx \n")
			datawriter.WriteString("    add rax, rbx \n")
			datawriter.WriteString("    push rax \n")
		case lexer.TOKEN_MINUS:
			datawriter.WriteString("    ;;-- minus %d -- \n")
			datawriter.WriteString("    pop rax \n")
			datawriter.WriteString("    pop rbx \n")
			datawriter.WriteString("    sub rbx, rax \n")
			datawriter.WriteString("    push rbx \n")
		case lexer.TOKEN_EQUALS:
			datawriter.WriteString("    ;;-- equals %d -- \n")
			datawriter.WriteString("    mov rcx, 0 \n")
			datawriter.WriteString("    mov rdx, 1 \n")
			datawriter.WriteString("    pop rax \n")
			datawriter.WriteString("    pop rbx \n")
			datawriter.WriteString("    cmp rax, rbx  \n")
			datawriter.WriteString("    cmove rcx,rdx  \n")
			datawriter.WriteString("    push rcx  \n")
		case lexer.TOKEN_IF:
			datawriter.WriteString("    ;;-- if %d -- \n")
			datawriter.WriteString("    pop rax \n")
			// PERFORM BITWISE AND AND SET VALUE TO ZERO FLAG
			datawriter.WriteString("    test rax, rax \n")
			// JUMP TO END ADDRESS GENERATED BY IF WHEN ZERO FLAG IS 1? ( so if rax is 0 )
			datawriter.WriteString(fmt.Sprintf("    jz addr_%d \n", token.Parameter))
		case lexer.TOKEN_ELSE:
			datawriter.WriteString("    ;;-- else -- \n")
			datawriter.WriteString(fmt.Sprintf("    jmp addr_%d\n", token.Parameter))
		case lexer.TOKEN_END:
			datawriter.WriteString("    ;;-- end %d -- \n")
			if token_index+1 != token.Parameter {
				datawriter.WriteString(fmt.Sprintf("    jmp addr_%d\n", token.Parameter))
			}
			continue
		case lexer.TOKEN_DUP:
			datawriter.WriteString("    ;;-- dup %d -- \n")
			datawriter.WriteString("    pop rax \n")
			datawriter.WriteString("    push rax \n")
			datawriter.WriteString("    push rax \n")
		case lexer.TOKEN_GREATER_THAN:
			datawriter.WriteString("    ;;-- greater than %d -- \n")
			datawriter.WriteString("    mov rcx, 0 \n")
			datawriter.WriteString("    mov rdx, 1 \n")
			datawriter.WriteString("    pop rbx \n")
			datawriter.WriteString("    pop rax \n")
			datawriter.WriteString("    cmp rax, rbx  \n")
			datawriter.WriteString("    cmovg rcx,rdx  \n")
			datawriter.WriteString("    push rcx  \n")
		case lexer.TOKEN_WHILE:
			datawriter.WriteString("    ;;-- while %d -- \n")
		case lexer.TOKEN_DO:
			datawriter.WriteString("    ;;-- do %d -- \n")
			datawriter.WriteString("    pop rax \n")
			datawriter.WriteString("    test rax, rax \n")
			datawriter.WriteString(fmt.Sprintf("    jz addr_%d \n", token.Parameter))
		case lexer.TOKEN_DUMP:
			datawriter.WriteString("    ;;-- dump %d -- \n")
			datawriter.WriteString("    pop rdi \n")
			datawriter.WriteString("    call dump\n")
		}
	}

	datawriter.WriteString("    mov rax, 60" + "\n")
	datawriter.WriteString("    mov rdi, 0" + "\n")
	datawriter.WriteString("    syscall" + "\n")
	datawriter.Flush()
	file.Close()

	utils.RunCMD("nasm -felf64 output.asm")
	utils.RunCMD("ld -o output output.o")
}

func (n *NAH) Interpret() {
	utils.CountTokensCheck(lexer.COUNT_TOKENS, 12, "./NAHI/NAHI.go", "Interpret")

	var programstack utils.Stack
	for t_token_index := 0; t_token_index < len(n.LEXER.Tokens); {
		token := n.LEXER.Tokens[t_token_index]

		switch token.TokenType {
		case lexer.TOKEN_PUSH:
			programstack.Push(token.Parameter)
		case lexer.TOKEN_PLUS:
			a, _ := programstack.Pop()
			b, _ := programstack.Pop()
			if reflect.TypeOf(a).Kind() == reflect.Int && reflect.TypeOf(b).Kind() == reflect.Int {
				a := a.(int)
				b := b.(int)
				programstack.Push(a + b)
			}
			// later on do string concat here maybe
		case lexer.TOKEN_MINUS:
			a, _ := programstack.Pop()
			b, _ := programstack.Pop()
			if reflect.TypeOf(a).Kind() == reflect.Int && reflect.TypeOf(b).Kind() == reflect.Int {
				a := a.(int)
				b := b.(int)
				programstack.Push(b - a)
			}
		case lexer.TOKEN_EQUALS:
			a, _ := programstack.Pop()
			b, _ := programstack.Pop()
			if reflect.TypeOf(a).Kind() == reflect.Int && reflect.TypeOf(b).Kind() == reflect.Int {
				a := a.(int)
				b := b.(int)
				if a == b {
					programstack.Push(1)
				} else {
					programstack.Push(0)
				}
			}
		case lexer.TOKEN_IF:
			a, _ := programstack.Pop()
			// if false jump to end
			if a == 0 {
				t_token_index = token.Parameter.(int)
				continue
			}
			// otherwise continue executing?
		case lexer.TOKEN_END:
			t_token_index = token.Parameter.(int)
			continue
		case lexer.TOKEN_ELSE:
			t_token_index = token.Parameter.(int)
			continue
		case lexer.TOKEN_DUMP:
			a, _ := programstack.Pop()
			fmt.Printf("%v \n", a)
		case lexer.TOKEN_DUP:
			a, _ := programstack.Pop()
			programstack.Push(a)
			programstack.Push(a)
		case lexer.TOKEN_GREATER_THAN:
			b, _ := programstack.Pop()
			a, _ := programstack.Pop()
			if a.(int) > b.(int) {
				programstack.Push(1)
			} else {
				programstack.Push(0)
			}
		case lexer.TOKEN_WHILE:
			t_token_index++
			continue
		case lexer.TOKEN_DO:
			a, _ := programstack.Pop()

			if a.(int) == 0 {
				t_token_index = token.Parameter.(int)
				continue
			} else {
				t_token_index++
				continue
			}
		default:
			fmt.Println("Unreachable")
		}
		t_token_index++
	}
}
