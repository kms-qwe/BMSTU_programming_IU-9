#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <string.h>
#include <stdbool.h>
#define STACK_SIZE (1 << 24)

typedef struct 
{
	long long *data;
	int cap;
	int top;
} STACK;

STACK* init_stack(int size)
{
	STACK* new_stack = malloc(sizeof(STACK));
	new_stack->data = calloc(size, sizeof(long long));
	new_stack->cap = size;
	new_stack->top = 0;
	return new_stack;
}

bool stack_empty(STACK *stack)
{
	return stack->top == 0;
}

void push(STACK *stack, long long x)
{
	if (stack->top == stack->cap)
		printf("ERROR: STACK OVERFLOW\n");
	stack->data[stack->top] = x;
	stack->top += 1;
}

long long pop(STACK *stack)
{
	if (stack_empty(stack))
		printf("ERROR: STACK UNDERFLOW\n");
	stack->top -= 1;
	return stack->data[stack->top];
}

long long max(long long a, long long b)
{
	return a > b ? a : b;
}

long long min(long long a, long long b)
{
	return a < b ? a : b;
}
int main(int argc, char const *argv[])
{
	STACK *stack = init_stack(STACK_SIZE);
	char *command = calloc(10, sizeof(char));
	long long number_on_stack;

	scanf("%s", command);
	while (strcmp(command, "END"))
	{
		if (!strcmp(command, "CONST"))
		{
			scanf("%lld", &number_on_stack);
			push(stack, number_on_stack);
		}

		if (!strcmp(command, "ADD"))
			push(stack, pop(stack) + pop(stack));

		if (!strcmp(command, "SUB"))
			push(stack, pop(stack) - pop(stack));

		if (!strcmp(command, "MUL"))
			push(stack, pop(stack) * pop(stack));

		if (!strcmp(command, "DIV"))
			push(stack, pop(stack) / pop(stack));

		if (!strcmp(command, "MAX"))
			push(stack, max(pop(stack), pop(stack)));

		if (!strcmp(command, "MIN"))
			push(stack, min(pop(stack), pop(stack)));

		if (!strcmp(command, "NEG"))
			push(stack, - pop(stack));

		if (!strcmp(command, "DUP"))
		{
			number_on_stack = pop(stack);
			push(stack, number_on_stack);
			push(stack, number_on_stack);
		}

		if (!strcmp(command, "SWAP"))
		{
			long long number_on_stack2 = pop(stack);
			number_on_stack = pop(stack);
			push(stack, number_on_stack2);
			push(stack, number_on_stack);
		}

		scanf("%s", command);
	} 
	
	printf("%lld\n", stack->data[stack->top - 1]);
	free(command);
	free(stack->data);
	free(stack);
	return 0;
}