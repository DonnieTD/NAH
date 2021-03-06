package nahi

import (
	"bufio"
	"fmt"
	"log"
	"os"

	lexer "github.com/DonnieTD/NAH/Lexer"
	utils "github.com/DonnieTD/NAH/Utils"
)

func (n *NAH) Compile() {
	utils.CountTokensCheck(lexer.COUNT_TOKENS, 26, "./NAHI/Compile.go:13", "Compile")

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
		case lexer.TOKEN_MEM:
			datawriter.WriteString("    ;;-- mem %d -- \n")
			datawriter.WriteString("    push mem \n")
		case lexer.TOKEN_STORE:
			datawriter.WriteString("    ;;-- store %d -- \n")
			datawriter.WriteString("    pop rbx \n")
			datawriter.WriteString("    pop rax \n")
			datawriter.WriteString("    mov [rax], bl \n")
		case lexer.TOKEN_LOAD:
			datawriter.WriteString("    ;;-- load %d -- \n")
			datawriter.WriteString("    pop rax \n")
			datawriter.WriteString("    xor rbx,rbx \n") //clean rbx
			datawriter.WriteString("    mov bl, [rax]\n")
			datawriter.WriteString("    push rbx \n")
		case lexer.TOKEN_SYSCALL3:
			datawriter.WriteString("    ;;-- syscall %d -- \n")
			datawriter.WriteString("    pop rax  \n")
			datawriter.WriteString("    pop rdi  \n")
			datawriter.WriteString("    pop rsi  \n")
			datawriter.WriteString("    pop rdx  \n")
			datawriter.WriteString("    syscall  \n")
		case lexer.TOKEN_SYSCALL1:
			datawriter.WriteString("    ;;-- syscall %d -- \n")
			datawriter.WriteString("    pop rax  \n")
			datawriter.WriteString("    pop rdi  \n")
			datawriter.WriteString("    syscall  \n")
		case lexer.TOKEN_DUMP:
			datawriter.WriteString("    ;;-- dump %d -- \n")
			datawriter.WriteString("    pop rdi \n")
			datawriter.WriteString("    call dump\n")
		case lexer.TOKEN_2DUP:
			datawriter.WriteString("    ;; -- 2dup -- \n")
			datawriter.WriteString("    pop rbx\n")
			datawriter.WriteString("    pop rax\n")
			datawriter.WriteString("    push rax\n")
			datawriter.WriteString("    push rbx\n")
			datawriter.WriteString("    push rax\n")
			datawriter.WriteString("    push rbx\n")
		case lexer.TOKEN_LESS_THAN:
			datawriter.WriteString("    ;; -- less than --\n")
			datawriter.WriteString("    mov rcx, 0\n")
			datawriter.WriteString("    mov rdx, 1\n")
			datawriter.WriteString("    pop rbx\n")
			datawriter.WriteString("    pop rax\n")
			datawriter.WriteString("    cmp rax, rbx\n")
			datawriter.WriteString("    cmovl rcx, rdx\n")
			datawriter.WriteString("    push rcx\n")
		case lexer.TOKEN_SWAP:
			datawriter.WriteString("    ;; -- swap --\n")
			datawriter.WriteString("    pop rax\n")
			datawriter.WriteString("    pop rbx\n")
			datawriter.WriteString("    push rax\n")
			datawriter.WriteString("    push rbx\n")
		case lexer.TOKEN_DROP:
			datawriter.WriteString("    ;; -- drop --\n")
			datawriter.WriteString("    pop rax\n")
		case lexer.TOKEN_SHL:
			datawriter.WriteString("    ;; -- shl --\n")
			datawriter.WriteString("    pop rcx\n")
			datawriter.WriteString("    pop rbx\n")
			datawriter.WriteString("    shl rbx, cl\n")
			datawriter.WriteString("    push rbx\n")
		case lexer.TOKEN_SHR:
			datawriter.WriteString("    ;; -- shr --\n")
			datawriter.WriteString("    pop rcx\n")
			datawriter.WriteString("    pop rbx\n")
			datawriter.WriteString("    shr rbx, cl\n")
			datawriter.WriteString("    push rbx\n")
		case lexer.TOKEN_BOR:
			datawriter.WriteString("    ;; -- bor --\n")
			datawriter.WriteString("    pop rax\n")
			datawriter.WriteString("    pop rbx\n")
			datawriter.WriteString("    or rbx, rax\n")
			datawriter.WriteString("    push rbx\n")
		case lexer.TOKEN_BAND:
			datawriter.WriteString("    ;; -- band --\n")
			datawriter.WriteString("    pop rax\n")
			datawriter.WriteString("    pop rbx\n")
			datawriter.WriteString("    and rbx, rax\n")
			datawriter.WriteString("    push rbx\n")
		case lexer.TOKEN_OVER:
			datawriter.WriteString("    ;; -- over --\n")
			datawriter.WriteString("    pop rax\n")
			datawriter.WriteString("    pop rbx\n")
			datawriter.WriteString("    push rbx\n")
			datawriter.WriteString("    push rax\n")
			datawriter.WriteString("    push rbx\n")
		}

	}

	datawriter.WriteString("    mov rax, 60 \n")
	datawriter.WriteString("    mov rdi, 0 \n")
	datawriter.WriteString("    syscall \n")

	// Memory segment
	datawriter.WriteString("segment .bss\n")
	datawriter.WriteString(fmt.Sprintf("mem: resb %d \n", MEM_CAPACITY))
	datawriter.Flush()
	file.Close()

	utils.RunCMD("nasm -felf64 output.asm")
	utils.RunCMD("ld -o output output.o")
}
