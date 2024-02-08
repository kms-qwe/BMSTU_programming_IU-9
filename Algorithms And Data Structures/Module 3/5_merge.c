#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <stdbool.h>
#include <string.h>

int compare(int A, int B)
{
	if(A < B) return 1;
	if(A > B) return -1;
	return 0;
}





void print_array(char *name, int *array, int len, int key)
{
	if(key)
	{
		// printf("________________\n");
		for(int i = 0; i < len; i++)
			printf("%s[%d] = %d\n", name, i, array[i]);
		// printf("________________\n");
	}
	else
	{
		for(int i = 0; i < len; i++)
			printf("%d\n", array[i]);
	}
}

typedef struct
{
	int index, key, low, cap;
	int *value;
} ELEMENT_QUEUE;

void swap(ELEMENT_QUEUE **a, ELEMENT_QUEUE **b)
{
	ELEMENT_QUEUE *tmp = *a;
	*a = *b;
	*b =  tmp;
}

ELEMENT_QUEUE* init_element_queue(int ind, int k, int cap, int *array)
{
	ELEMENT_QUEUE *new_element = malloc(sizeof(ELEMENT_QUEUE));
	new_element->index = ind;
	new_element->key = k;
	new_element->value = array;
	new_element->low = 0;
	new_element->cap = cap;
	return new_element;
}

typedef struct
{
	ELEMENT_QUEUE **heap;
	int cap, count;
} PRIORITY_QUEUE;

PRIORITY_QUEUE* init_priority_queue(int size)
{
	PRIORITY_QUEUE *queue = malloc(sizeof(PRIORITY_QUEUE));
	queue->heap = calloc(size, sizeof(ELEMENT_QUEUE*));
	for(int i = 0; i < size; i++)
		queue->heap[i] = NULL;
	queue->cap = size;
	queue->count = 0;
	return queue;
}

bool priority_queue_empty(PRIORITY_QUEUE *queue)
{
	return queue->count == 0;
}

ELEMENT_QUEUE* max_priority_queue(PRIORITY_QUEUE *queue)
{
	if(priority_queue_empty(queue)) 
		printf("ERROR: PRIORITY_QUEUE IS EMPTY (max_priority_queue)\n");
	return queue->heap[0];
}

void print_queue_element(ELEMENT_QUEUE *element)
{
	printf("\n");
	printf("index = %d, key = %d, low = %d, cap = %d\n", 
		element->index, element->key, element->low, element->cap);
	print_array("q_el", element->value, element->cap, 1);
	// printf("_________________\n");	
}

void print_queue(PRIORITY_QUEUE *queue)
{
	printf("-----------------\n");
	printf("cap = %d, count = %d\n", queue->cap, queue->count);
	for(int i = 0; i < queue->count; i++)
		if(queue->heap[i] == NULL)
			printf("queue->heap[%d] == NULL\n", i);
		else
			print_queue_element(queue->heap[i]);
	printf("-----------------\n");	
}

void priority_queue_insert(PRIORITY_QUEUE *queue, ELEMENT_QUEUE* new_element)
{
	if(queue->count == queue->cap)
		printf("ERROR: PRIORITY_QUEUE OVERFLOW\n");
	int i = queue->count;
	queue->count++;
	queue->heap[i] = new_element;
	while(i > 0 and compare(queue->heap[(i - 1)/2]->key, queue->heap[i]->key) < 0)
	{
		swap(queue->heap + (i - 1)/2, queue->heap + i);
		queue->heap[i]->index = i;
		i = (i - 1)/2;		
	}
	queue->heap[i]->index = i;
}



void heapify(int root_ind, int len, ELEMENT_QUEUE **heap)
{
	while(true)
	{
		int init_root_ind = root_ind;
		int left_child_ind = 2 * root_ind + 1;
		int right_child_ind = left_child_ind + 1;
		if(left_child_ind < len and compare(heap[root_ind]->key, heap[left_child_ind]->key) < 0)
			root_ind = left_child_ind;
		if(right_child_ind < len and compare(heap[root_ind]->key, heap[right_child_ind]->key) < 0)
			root_ind = right_child_ind;
		if(init_root_ind == root_ind)
			break;
		swap(heap + root_ind, heap + init_root_ind);
		heap[root_ind]->index = root_ind;
		heap[init_root_ind]->index = init_root_ind;
	}
}

ELEMENT_QUEUE* extract_max(PRIORITY_QUEUE *queue)
{
	if(priority_queue_empty(queue))
		printf("ERROR: PRIORITY_QUEUE IS EMPTY (extract_max)\n");
	ELEMENT_QUEUE *ptr_on_max = queue->heap[0];
	queue->count--;
	if(!priority_queue_empty(queue))
	{
		queue->heap[0] = queue->heap[queue->count];
		queue->heap[0]->index = 0;
		heapify(0, queue->count, queue->heap);
	}
	return ptr_on_max;
}

void increase_key(PRIORITY_QUEUE *queue, ELEMENT_QUEUE *ptr, int new_key)
{
	int new_index = ptr->index;
	ptr->key = new_key;
	while(new_index > 0 and compare(queue->heap[(new_index - 1)/2]->key, new_key) < 0)
	{
		swap(queue->heap + (new_index - 1)/2, queue->heap + new_index);
		queue->heap[new_index]->index = new_index;
		new_index = (new_index - 1)/2;
	}
	ptr->index = new_index;
}

void decrease_key(PRIORITY_QUEUE *queue, ELEMENT_QUEUE *ptr, int new_key)
{
	ptr->key = new_key;
	heapify(ptr->index, queue->count, queue->heap);
}

int* merge(int **array, int len_array_of_arrays, int *len_arrays, int len_merge)
{

	int top_merge_array = 0;
	int *merge_array = malloc(len_merge * sizeof(int));
	PRIORITY_QUEUE *queue = init_priority_queue(len_array_of_arrays);

	for(int i = 0; i < len_array_of_arrays; i++)
	{
		ELEMENT_QUEUE *new_element = init_element_queue(-1, array[i][0], len_arrays[i] ,array[i]);
		priority_queue_insert(queue, new_element);
	}

	while(!priority_queue_empty(queue))
	{
		ELEMENT_QUEUE *maximum = max_priority_queue(queue);
		merge_array[top_merge_array] = maximum->value[maximum->low];
		top_merge_array++;
		maximum->low++;
		if(maximum->low == maximum->cap)
		{
			maximum = extract_max(queue);
			free(maximum);
		}
		else
		{
			decrease_key(queue, maximum, maximum->value[maximum->low]);
		}
	}

	free(queue->heap);
	free(queue);
	return merge_array;
}

int main()
{

	int len_array_of_arrays, len_merge = 0;
	scanf("%d", &len_array_of_arrays);
	int **array = calloc(len_array_of_arrays, sizeof(int*));
	int *len_arrays = calloc(len_array_of_arrays, sizeof(int));


	for(int i = 0; i < len_array_of_arrays; i++)
	{
		scanf("%d", &len_arrays[i]);
		len_merge += len_arrays[i];
	}

	for(int i = 0; i < len_array_of_arrays; i++)
	{
		array[i] = calloc(len_arrays[i], sizeof(int));
		for(int j = 0; j < len_arrays[i]; j++)
			scanf("%d", &array[i][j]);
	}

	int *merge_array = merge(array, len_array_of_arrays, len_arrays, len_merge);
	print_array("merge_array", merge_array, len_merge, 0);

	for(int i = 0; i < len_array_of_arrays; i++)
		free(array[i]);
	free(array);
	free(len_arrays);
	free(merge_array);
	return 0;
}
