#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <stdbool.h>
#include <string.h>

int max(int a, int b)
{
	return a > b ? a : b;
}

int compare(int A, int B)
{
	if(A < B) return 1;
	if(A > B) return -1;
	return 0;
}

// void swap(void *a, void *b, size_t size)
// {

// 	void *tmp = malloc(size);
// 	memcpy(tmp, a, size);
// 	memcpy(a, b, size);
// 	memcpy(b, tmp, size);
// 	free(tmp);
// }

void swap(int *a, int *b, size_t size)
{
	*a = (*a)^(*b);
	*b = (*a)^(*b);
	*a = (*b)^(*a);
}

typedef struct
{
	int *heap;
	int cap, count;
} PRIORITY_QUEUE;

PRIORITY_QUEUE* init_priority_queue(int size)
{
	PRIORITY_QUEUE *queue = malloc(sizeof(PRIORITY_QUEUE));
	queue->heap = calloc(size, sizeof(int));
	queue->cap = size;
	queue->count = 0;
	return queue;
}

bool priority_queue_empty(PRIORITY_QUEUE *queue)
{
	return queue->count == 0;
}

int max_priority_queue(PRIORITY_QUEUE *queue)
{
	if(priority_queue_empty(queue)) 
		printf("ERROR: PRIORITY_QUEUE IS EMPTY (max_priority_queue)\n");
	return queue->heap[0];
}

void priority_queue_insert(PRIORITY_QUEUE *queue, int new_element)
{
	if(queue->count == queue->cap)
		printf("ERROR: PRIORITY_QUEUE OVERFLOW\n");
	int i = queue->count;
	queue->count++;
	queue->heap[i] = new_element;
	while(i > 0 and compare(queue->heap[(i - 1)/2], queue->heap[i]) < 0)
	{
		swap(queue->heap + (i - 1)/2, queue->heap + i, sizeof(int));
		i = (i - 1)/2;		
	}
}



void heapify(int root_ind, int len, int *heap)
{
	while(true)
	{
		int init_root_ind = root_ind;
		int left_child_ind = 2 * root_ind + 1;
		int right_child_ind = left_child_ind + 1;
		if(left_child_ind < len and compare(heap[root_ind], heap[left_child_ind]) < 0)
			root_ind = left_child_ind;
		if(right_child_ind < len and compare(heap[root_ind], heap[right_child_ind]) < 0)
			root_ind = right_child_ind;
		if(init_root_ind == root_ind)
			break;
		swap(heap + root_ind, heap + init_root_ind, sizeof(int));
	}
}

int extract_max(PRIORITY_QUEUE *queue)
{
	if(priority_queue_empty(queue))
		printf("ERROR: PRIORITY_QUEUE IS EMPTY (extract_max)\n");
	int max_element = queue->heap[0];
	queue->count--;
	if(!priority_queue_empty(queue))
	{
		queue->heap[0] = queue->heap[queue->count];
		heapify(0, queue->count, queue->heap);
	}
	return max_element;
}

void increase_key(PRIORITY_QUEUE *queue, int index, int new_value)
{
	queue->heap[index] = new_value;
	while(index > 0 and compare(queue->heap[(index - 1)/2], new_value) < 0)
	{
		swap(queue->heap + (index - 1)/2, queue->heap + index, sizeof(int));
		index = (index - 1)/2;
	}
}

void decrease_key(PRIORITY_QUEUE *queue, int index, int new_value)
{
	queue->heap[index] = new_value;
	heapify(index, queue->count, queue->heap);
}

int min_time_cluster()
{
	int N, M, t1, t2;
	scanf("%d %d", &N, &M);
	PRIORITY_QUEUE *nodes = init_priority_queue(N);

	for(int i = 0; i< N; i++)
		priority_queue_insert(nodes, 0);
		 
	for(int i = 0; i < M; i++)
	{
		scanf("%d %d", &t1, &t2);
		int current_min_time = max_priority_queue(nodes);
		current_min_time = max(current_min_time, t1) + t2;
		decrease_key(nodes, 0, current_min_time);	
	}
	int min_time;
	while(!priority_queue_empty(nodes))
		min_time = extract_max(nodes);

	free(nodes->heap);
	free(nodes);
	return min_time;
}
int main()
{	
	int min_time = min_time_cluster();
	printf("%d\n", min_time);
	return 0;
}
