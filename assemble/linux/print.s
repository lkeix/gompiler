
.data
.Hello:
  .string "hello world\n"

.text

.global _start
_start:
  callq main.main

simple.print:
  movq $2, %rdi # rdi is set to 2 as stderr
  movq 16(%rsp), %rsi # rsi is set to the string to print
  movq 8(%rsp),  %rdx # rdx is set to the length of the string
  movq $1, %rax # rax is set to 1 as stdout
  syscall
  ret

os.Exit:
  movq 8(%rsp), %rdi
  movq $60, %rax
  movq $0, %rdi
  syscall

main.main:
  leaq .Hello, %rax # the reason using leaq is .Hello variable is stored in static address
  pushq %rax # push the address of the string to the stack
  pushq $12 # string length is stored in stack
  call simple.print
  popq %rax
  call os.Exit
