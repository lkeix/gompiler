  .globl main

main:
  push    %rbp
  movq    %rsp, %rbp
  movl    $0,  %eax
  pop     %rbp
  ret
