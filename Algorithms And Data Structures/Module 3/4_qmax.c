#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <string.h>
#include <stdbool.h>
#include <limits.h>
#define SIZE (1 << 25)
void print_bin_long(long long x)
{
	printf("__________\n");
	printf("x = %lld\n", x);
	for(int i = 63; i >= 0; i--)
		printf("%lld", (x & ((long long)1 << i)) >> i);
	printf("\n");
}

long long write_left(long long orig, int left)
{
	long long long_left = 0;
	long long min_mask = -1 & (((long long)1 << 32) - 1);
	for(int i = 0; i < 32; i++)
		long_left |= ((left & (1 << i)) >= 0 ? 
			(left & (1 << i)) : (long long)1 << 31);
	return (orig & min_mask) | (long_left << 32);
}

long long write_right(long long orig, int right)
{
	long long long_right = 0;
	long long min_mask = -1 & (((long long)1 << 32) - 1);
	for(int i = 0; i < 32; i++)
		long_right |= ((right & (1 << i)) >= 0 ? 
			(right & (1 << i)) : (long long)1 << 31);	
	return (orig  & (~min_mask)) | long_right;
}

int get_left(long long x)
{
	long long min_mask = -1 & (((long long)1 << 32) - 1);
	return (int)((x & (~min_mask)) >> 32);
}

int get_right(long long x)
{
	long long min_mask = -1 & (((long long)1 << 32) - 1);
	return (int)(x & min_mask);
}

int max(int x, int y)
{
	return x > y ? x : y;
}

typedef struct 
{
	long long *data;
	int cap;
	int top1;
	int top2;
} DOUBLE_STACK;

DOUBLE_STACK* init_double_stack(int size)
{
	DOUBLE_STACK *double_stack = malloc(sizeof(DOUBLE_STACK));
	double_stack->data = calloc(size, sizeof(long long));
	double_stack->cap = size;
	double_stack->top1 = 0;
	double_stack->top2 = size - 1;
	return double_stack;
}

bool stack_empty1(DOUBLE_STACK *double_stack)
{
	return double_stack->top1 == 0;
}

bool stack_empty2(DOUBLE_STACK *double_stack)
{
	return double_stack->top2 == double_stack->cap - 1;
}

void push1(DOUBLE_STACK *double_stack, int new_element)
{
	if(double_stack->top2 < double_stack->top1)
		printf("ERROR: DOUBLE STACK OVERFLOW (push1)\n");

	double_stack->data[double_stack->top1] = 
		write_right(double_stack->data[double_stack->top1], new_element);	
	if(stack_empty1(double_stack))
		double_stack->data[double_stack->top1] = 
			write_left(double_stack->data[double_stack->top1], new_element);
	else
		double_stack->data[double_stack->top1] = 
			write_left(double_stack->data[double_stack->top1], 
				max(new_element, get_left(double_stack->data[double_stack->top1 - 1])));	
	double_stack->top1++;
}

void push2(DOUBLE_STACK *double_stack, int new_element)
{
	if(double_stack->top2 < double_stack->top1)
		printf("ERROR: DOUBLE STACK OVERFLOW (push2)\n");
	double_stack->data[double_stack->top2] = 
		write_right(double_stack->data[double_stack->top2], new_element);
	if(stack_empty2(double_stack))
		double_stack->data[double_stack->top2] = 
			write_left(double_stack->data[double_stack->top2], new_element);
	else
		double_stack->data[double_stack->top2] = 
			write_left(double_stack->data[double_stack->top2], 
				max(new_element, get_left(double_stack->data[double_stack->top2 + 1])));
	double_stack->top2--;
}

int pop1(DOUBLE_STACK *double_stack)
{
	if(stack_empty1(double_stack))
		printf("ERROR: DOUBLE STACK UNDERFLOW (pop1)\n");
	double_stack->top1--;
	return double_stack->data[double_stack->top1];
}

int pop2(DOUBLE_STACK *double_stack)
{
	if(stack_empty2(double_stack))
		printf("ERROR: DOUBLE STACK UNDERFLOW (pop2)\n");
	double_stack->top2++;
	return double_stack->data[double_stack->top2];
}

int get_max_1(DOUBLE_STACK *double_stack)
{
	return stack_empty1(double_stack) ? INT_MIN : get_left(double_stack->data[double_stack->top1 - 1]);
}

int get_max_2(DOUBLE_STACK *double_stack)
{
	return stack_empty2(double_stack) ? INT_MIN : get_left(double_stack->data[double_stack->top2 + 1]);
}

int get_max(DOUBLE_STACK *double_stack)
{
	return max(get_max_1(double_stack), get_max_2(double_stack));
}

void free_double_stack(DOUBLE_STACK *double_stack)
{
	free(double_stack->data);
	free(double_stack);
}

DOUBLE_STACK* init_queue_on_stack(int size)
{
	return init_double_stack(size);
}

bool queue_empty(DOUBLE_STACK *queue)
{
	return stack_empty1(queue) and stack_empty2(queue);
}

void enqueue(DOUBLE_STACK *queue, int new_element)
{
	push1(queue, new_element);
}



int dequeue(DOUBLE_STACK *queue)
{
	if(stack_empty2(queue))
		while(!stack_empty1(queue))
			push2(queue, pop1(queue));

	return get_right(pop2(queue));
}


void print_queue(DOUBLE_STACK *queue)
{
	printf("__________\n");
	printf("cap = %d, top1 = %d, top2 = %d, max = %d\n",
		queue->cap, queue->top1, queue->top2, get_max(queue));
	for(int i = 0; i < queue->cap; i++)
	{
		printf("queue[%d] = [%d, %d]\n", 
			i, get_left(queue->data[i]), get_right(queue->data[i]));
	}
	printf("__________\n");

}

int max_in_queue(DOUBLE_STACK *queue)
{
	return get_max(queue);
}

void free_queue(DOUBLE_STACK *queue)
{
	free_double_stack(queue);
}

void qmax()
{
	DOUBLE_STACK *queue = init_queue_on_stack(SIZE);
	int number_in_queue;
	char *command = calloc(10, sizeof(char));
	scanf("%s", command);

	while(strcmp(command, "END"))
	{
		if(!strcmp(command, "ENQ"))
		{
			scanf("%d", &number_in_queue);
			enqueue(queue, number_in_queue);
		}

		if(!strcmp(command, "DEQ"))
			printf("%d\n", dequeue(queue));

		if(!strcmp(command, "MAX"))
			printf("%d\n", max_in_queue(queue));

		if(!strcmp(command, "EMPTY"))
			printf("%s\n", queue_empty(queue) ? "true" : "false");
		scanf("%s", command);
	}

	free(command);
	free_queue(queue);
}

int main(int argc, char const *argv[])
{
	qmax();
	return 0;
}