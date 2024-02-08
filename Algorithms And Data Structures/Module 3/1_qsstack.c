#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <stdbool.h>
#include <string.h>

void print_array(char *name, long long *array, int len, int key)
{
	if(key)
	{
		printf("________________\n");
		for(int i = 0; i < len; i++)
			printf("%s[%d] = %lld\n", name, i, array[i]);
		printf("________________\n");
	}
	else
	{
		for(int i = 0; i < len; i++)
			printf("%lld\n", array[i]);
	}
}

typedef struct 
{
	int low, high;
} TASK;

TASK* init_task(int low, int high)
{
	TASK *task = malloc(sizeof(TASK));
	task->low = low;
	task->high = high;
	return task;
}

void free_task(TASK *task)
{
	free(task);
}

typedef struct
{
	TASK **data;
	int cap, top;
} STACK_OF_TASK;

STACK_OF_TASK* init_stack(int size)
{
	STACK_OF_TASK *stack = malloc(sizeof(STACK_OF_TASK));
	stack->data = calloc(size, sizeof(TASK*));
	stack->cap = size;
	stack->top = 0;
	return stack;
}

void free_stack_of_task(STACK_OF_TASK *stack_of_task)
{
	for(int i = 0; i < stack_of_task->top; i++)
		free_task(stack_of_task->data[i]);
	free(stack_of_task->data);
	free(stack_of_task);
}

bool stack_empty(STACK_OF_TASK *stack)
{
	return stack->top == 0;
}

void push(STACK_OF_TASK *stack, TASK* new_task)
{
	if(stack->top == stack->cap)
		printf("ERROR: STACK OVERFLOW\n");
	stack->data[stack->top] = new_task;
	stack->top++;
}

TASK* pop(STACK_OF_TASK *stack)
{
	if(stack_empty(stack))
		printf("ERROR: STACK UNDERFLOW\n");
	stack->top--;
	return stack->data[stack->top];
}

long long compare(const void *a, const void *b)
{
	return *(long long *)a - *(long long *)b;
}

void swap(void *a, void *b, size_t size)
{
	void *tmp = malloc(size);
	memcpy(tmp, b, size);
	memcpy(b, a, size);
	memcpy(a, tmp, size);
	free(tmp);
}

void print_stack_of_task(char *name, STACK_OF_TASK *stack_of_task)
{
	printf("________________\n");
	printf("cap = %d, top = %d\n", stack_of_task->cap, stack_of_task->top);
	for(int i = 0; i < stack_of_task->top; i++)
		printf("%s[%d] = [%d, %d]\n", name, i, (stack_of_task->data[i])->low, (stack_of_task->data[i])->high);
	printf("________________\n");
}

int partition(long long(*compare)(const void *a, const void *b), 
	void(*swap)(void *a, void *b, size_t size), int low, int high, long long *array)
{
	int i = low;
	for(int j = low; j < high; j++)
		if(compare(array + j, array + high) < 0)
		{
			swap(array + i, array + j, sizeof(long long));
			i++;
		}
	swap(array + i, array + high, sizeof(long long));
	return i;
}

void quick_sort_rec(long long(*compare)(const void *a, const void *b), 
	void(*swap)(void *a, void *b, size_t size), STACK_OF_TASK *stack_of_task, long long *array)
{
	while(!stack_empty(stack_of_task))
	{
		TASK *current_task = pop(stack_of_task);
		if(current_task->low < current_task->high)
		{
			int q = partition(compare, swap, current_task->low, current_task->high, array);
			TASK *task_low_to_q = init_task(current_task->low, q - 1);
			TASK *task_q_to_high = init_task(q + 1, current_task->high);
			push(stack_of_task, task_low_to_q);
			push(stack_of_task, task_q_to_high);
		}


		free_task(current_task);
	}	
}

void quick_sort(long long(*compare)(const void *a, const void *b), void(*swap)(void *a, void *b, size_t size), long long *array, int len)
{
	STACK_OF_TASK *stack_of_task = init_stack(len);
	TASK *first_task = init_task(0, len - 1);
	push(stack_of_task, first_task);
	quick_sort_rec(compare,swap, stack_of_task, array);
	if(!stack_empty(stack_of_task))
		printf("ERROR: STACK_OF_TASK IS NOT EMPTY AFTER SORT\n");
	free_stack_of_task(stack_of_task);
}





int main(int argc, char const *argv[])
{
	int len;
	scanf("%d", &len);
	long long *array = malloc(len * sizeof(long long));
	for(int i = 0; i < len; i++)
		scanf("%lld", array + i);
	quick_sort(compare, swap, array, len);
	print_array("Array", array, len, 0);
	free(array);
	return 0;
}