#!/usr/bin/env python
from enum import IntEnum
import struct

STACK_SIZE_MAX = 2048
class OpCode(IntEnum):

    # stack operation
    PUSH_INT           =  1,
    PUSH_DOUBLE        =  2,
    PUSH_STRING        =  3,
    PUSH_NULL          =  4,
    POP_INT            =  5,
    POP_DOUBLE         =  6,
    POP_STRING         =  7,
    POP_NULL           =  8,

    # calc operation
    ADD_INT            =  9,
    ADD_DOUBLE         = 10,
    SUB_INT            = 11,
    SUB_DOUBLE         = 12,
    MUL_INT            = 13,
    MUL_DOUBLE         = 14,
    DIV_INT            = 15,
    DIV_DOUBLE         = 16,
    MOD_INT            = 17,
    MOD_DOUBLE         = 18,

    CAST_INT_TO_DOUBLE = 19,
    CAST_DOUBLE_TO_INT = 20,

    # logic operation
    EQ_INT             = 21,
    EQ_DOUBLE          = 22,
    GT_INT             = 23,
    GT_DOUBLE          = 24,
    GE_INT             = 25,
    GE_DOUBLE          = 26,

    LT_INT             = 27,
    LT_DOUBLE          = 28,
    LE_INT             = 29,
    LE_DOUBLE          = 30,
    NE_INT             = 31,
    NE_DOUBLE          = 32,

    LOGICAL_AND        = 33,
    LOGICAL_OR         = 34,
    LOGICAL_NOT        = 35,

    # control operation
    JMP                = 36,
    JMP_IF_TRUE        = 37,
    JMP_IF_FALSE       = 38

class OpCodeInfo(object):
    def __init__(self,nemonic, param, stack_increment):
        self.nemonic         = nemonic
        self.param           = param
        self.stack_increment = stack_increment

class SchemeVm:
    def __init__(self):
        self.accumelator = 0
        self.pc = 0
        self.stack = [0] * STACK_SIZE_MAX
        self.stack_pos = STACK_SIZE_MAX - 1
        self.stack_ebp = self.stack_pos

    def int_stack_write(self, ivalue):
        stack = self.stack

        stack[ self.stack_pos ] = struct.unpack("i", ivalue)[0]
        print(stack[ self.stack_pos ])
        self.stack_pos -= 1

        self.pc = self.pc + 1 + 4

    def int_add(self):
        stack = self.stack

        stack[ self.stack_pos ] = 0
        for i in (stack[self.stack_pos + 1 : self.stack_ebp + 1]):
            stack[ self.stack_pos ] += i

        self.accumelator = self.stack[ self.stack_pos ]
        self.stack_ebp = self.stack_pos
        self.stack_pos -= self.stack_pos

        self.pc = self.pc + 1
        
    def print_accumelator(self):
        print(self.accumelator)

def createInfo():
    asmInfo = [OpCodeInfo("dummy", "", 0),
               OpCodeInfo("push_int", "p", 1),
               OpCodeInfo("push_double", "p", 1),
               OpCodeInfo("push_string", "p", 1),
               OpCodeInfo("push_null", "p", 1),
               OpCodeInfo("pop_int", "p", 1),
               OpCodeInfo("pop_double", "p", 1),
               OpCodeInfo("pop_string", "p", 1),
               OpCodeInfo("pop_null","",1),
               OpCodeInfo("add_int","",-1),
               OpCodeInfo("add_double","",-1),
               OpCodeInfo("sub_int","",-1),
               OpCodeInfo("sub_double","",-1),
               OpCodeInfo("mul_int","",-1),
               OpCodeInfo("mul_double","",-1),
               OpCodeInfo("div_int","",-1),
               OpCodeInfo("div_double","",-1),
               OpCodeInfo("mod_int","",-1),
               OpCodeInfo("mod_double","",-1),
               OpCodeInfo("cast_int_to_double", "", 0),
               OpCodeInfo("cast_double_to_int", "", 0),
               OpCodeInfo("eq_int", "", -1),
               OpCodeInfo("eq_double", "", -1),
               OpCodeInfo("gt_int", "", -1),
               OpCodeInfo("gt_double", "", -1),
               OpCodeInfo("ge_int", "", -1),
               OpCodeInfo("ge_double", "", -1),
               OpCodeInfo("lt_int", "", -1),
               OpCodeInfo("lt_double", "", -1),
               OpCodeInfo("le_int", "", -1),
               OpCodeInfo("le_double", "", -1),
               OpCodeInfo("ne_int", "", -1),
               OpCodeInfo("ne_double", "", -1),
               OpCodeInfo("logical_and", "", -1),
               OpCodeInfo("logical_or", "", -1),
               OpCodeInfo("logical_not", "", -1),
               OpCodeInfo("jmp", "", 0),
               OpCodeInfo("jmp_if_true", "", -1),
               OpCodeInfo("jmp_if_false", "", -1)]
    return asmInfo

def execute(vm, byte_code,info):
    code_size = len(byte_code)

    while (vm.pc < code_size):
        if (byte_code[vm.pc] == OpCode.PUSH_INT):
            #print(info[byte_code[vm.pc]].nemonic)
            vm.int_stack_write(byte_code[vm.pc + 1 : vm.pc + 1 + 4])


        if (byte_code[vm.pc] == OpCode.ADD_INT):
            #print(info[byte_code[vm.pc]].nemonic)
            vm.int_add()
            vm.print_accumelator()

# Main
execute(SchemeVm(), b'\x01\x0a\x00\x00\x00\x01\x0b\x00\x00\x00\x09', createInfo())
