#undef _FORTIFY_SOURCE

#include <cstdint>

#include <llvm/ADT/APInt.h>
#include <llvm/IR/BasicBlock.h>
#include <llvm/IR/Constants.h>
#include <llvm/IR/Function.h>
#include <llvm/IR/IRBuilder.h>
#include <llvm/IR/LLVMContext.h>
#include <llvm/IR/Module.h>
#include <llvm/IR/Type.h>
#include <llvm/IR/Verifier.h>
#include <llvm/IR/NoFolder.h>

int main() {
    llvm::LLVMContext context;

    llvm::Module module("simple module", context);

    llvm::IRBuilder<llvm::NoFolder> builder(context, llvm::NoFolder());

    llvm::FunctionType* funcType =
        llvm::FunctionType::get(llvm::Type::getInt32Ty(context), false);

    llvm::Function* mainFunc =
        llvm::Function::Create(funcType, llvm::Function::ExternalLinkage, "main", module);

    llvm::BasicBlock* entryBlock =
        llvm::BasicBlock::Create(context, "entry", mainFunc);

    builder.SetInsertPoint(entryBlock);

    llvm::Value* lhs =
        llvm::ConstantInt::get(context, llvm::APInt(32, 353));

    llvm::Value* rhs =
        llvm::ConstantInt::get(context, llvm::APInt(32, 48));

    llvm::Value* result = builder.CreateAdd(lhs, rhs, "sum");

    builder.CreateRet(result);

    llvm::verifyFunction(*mainFunc);

    module.print(llvm::errs(), nullptr);

    return 0;
}
