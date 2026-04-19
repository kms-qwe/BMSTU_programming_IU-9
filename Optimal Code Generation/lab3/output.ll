; ModuleID = 'mod'
source_filename = "mod"

define i32 @main() {
entry:
  %a = alloca i32, align 4
  store i32 12, ptr %a, align 4
  %b = alloca i32, align 4
  store i32 4, ptr %b, align 4
  %a1 = load i32, ptr %a, align 4
  %b2 = load i32, ptr %b, align 4
  %sub = sub i32 %a1, %b2
  %result = alloca i32, align 4
  store i32 %sub, ptr %result, align 4
  %result3 = load i32, ptr %result, align 4
  %gt = icmp sgt i32 %result3, 6
  br i1 %gt, label %then, label %merge

then:                                             ; preds = %entry
  %result4 = load i32, ptr %result, align 4
  %add = add i32 %result4, 2
  store i32 %add, ptr %result, align 4
  br label %merge

merge:                                            ; preds = %then, %entry
  br label %cond

cond:                                             ; preds = %loop, %merge
  %b5 = load i32, ptr %b, align 4
  %gt6 = icmp sgt i32 %b5, 0
  br i1 %gt6, label %loop, label %after

loop:                                             ; preds = %cond
  %result7 = load i32, ptr %result, align 4
  %b8 = load i32, ptr %b, align 4
  %add9 = add i32 %result7, %b8
  store i32 %add9, ptr %result, align 4
  %b10 = load i32, ptr %b, align 4
  %sub11 = sub i32 %b10, 1
  store i32 %sub11, ptr %b, align 4
  br label %cond

after:                                            ; preds = %cond
  %result12 = load i32, ptr %result, align 4
  ret i32 %result12
}
