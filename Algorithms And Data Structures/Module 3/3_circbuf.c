#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <string.h>
#include <stdbool.h>
#define INIT_SIZE 4
typedef struct
{
	int *data;
	int cap;
	int count;
	int head;
	int tail;
} QUEUE;

QUEUE* init_queue(int size)
{
	QUEUE *queue = malloc(sizeof(QUEUE));
	queue->data = calloc(size, sizeof(int));
	queue->cap = size;
	queue->count = queue->head = queue->tail = 0;
	return queue;
}

void print_queue(QUEUE *queue)
{
	printf("____________\n");
	printf("cap = %d, count = %d, head = %d, tail = %d\n", 
		queue->cap, queue->count, queue->head, queue->tail);
	for(int i = 0; i < queue->cap; i++)
		printf("queue[%d] = %d\n", i, queue->data[i]);
	printf("____________\n");
}

void free_queue(QUEUE *queue)
{
	free(queue->data);
	free(queue);
}

bool queue_empty(QUEUE *queue)
{
	return queue->count == 0;
}

bool queue_full(QUEUE *queue)
{
	return queue->count == queue->cap;
}

void enqueue(QUEUE *queue, int x)
{
	if(queue_full(queue))
		printf("ERROR: QUEUE OVERFLOW\n");
	queue->data[queue->tail] = x;
	queue->tail = (queue->tail + 1) % queue->cap;
	queue->count++;
}

int dequeue(QUEUE *queue)
{
	if(queue_empty(queue))
		printf("ERROR: QUEUE UNDERFLOW\n");
	int x = queue->data[queue->head];
	queue->head = (queue->head + 1) % queue->cap;
	queue->count--;
	return x;
}



QUEUE* double_queue(QUEUE *queue)
{
	QUEUE *new_queue = init_queue(queue->cap*2);
	for(int i = 0; i < queue->cap; i++)
		new_queue->data[i] = queue->data[(queue->head + i) % queue->cap];
	new_queue->count = queue->count;
	new_queue->tail = queue->cap;
	free_queue(queue);
	return new_queue;
}


void circbuf()
{
	QUEUE *queue = init_queue(INIT_SIZE);
	char* command = calloc(10, sizeof(char));
	int number_in_queue;
	scanf("%s", command);

	while(strcmp(command, "END"))
	{

		if(!strcmp(command, "ENQ"))
		{
			if(queue_full(queue))
				queue = double_queue(queue);
			scanf("%d", &number_in_queue);
			enqueue(queue, number_in_queue);
		}

		if(!strcmp(command, "DEQ"))
			printf("%d\n", dequeue(queue));

		if(!strcmp(command, "EMPTY"))
			printf("%s\n", queue_empty(queue) ? "true" : "false");
		// print_queue(queue);
		scanf("%s", command);
	}

	free(command);
	free_queue(queue);
}

int main(int argc, char const *argv[])
{
	circbuf();
	return 0;
}
