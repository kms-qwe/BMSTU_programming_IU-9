#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <stdbool.h>
#include <string.h>
#include <time.h>
#include <limits.h>

typedef struct ELEM
{
		struct ELEM *prev, *next;
		int value;
} ELEM;

ELEM* init_double_linked_list()
{
	ELEM* list = malloc(sizeof(ELEM));
	list->prev = list;
	list->next = list;
	list->value = INT_MAX;
	return list;
}

void print_list(ELEM* list, int key)
{
	if(key) 
	{
		printf("RESTRICTIVE ELEMENT\n");
		ELEM *step = list->next;
		int index = 0;
		while(step != list)
		{
			printf("list[%d] = %d\n", index, step->value);
			step = step->next;
			index++;
		}
	}
	else
	{
		ELEM *step = list->next;
		int index = 0;
		while(step != list)
		{
			printf("%d\n", step->value);
			step = step->next;
			index++;
		}
	}
}

bool list_empty(ELEM *list)
{
	return list->next == list;
}

int list_length(ELEM *list)
{
	int len = 0;
	ELEM *step = list;
	while(step->next != list)
	{
		len++;
		step = step->next;
	}
	return len;
}

void insert_after(ELEM* after_this, ELEM* insert_this)
{
	ELEM* after_after = after_this->next;
	after_this->next = insert_this;
	insert_this->prev = after_this;
	insert_this->next = after_after;
	after_after->prev = insert_this;
}

void insert_before(ELEM* before_this, ELEM* insert_this)
{
	ELEM* before_before = before_this->prev;
	before_this->prev = insert_this;
	insert_this->next = before_this;
	insert_this->prev = before_before;
	before_before->next = insert_this;
}

void insert_before_head(ELEM *list, ELEM *new_element)
{
	insert_after(list, new_element);
}

void insert_after_tail(ELEM *list, ELEM *new_element)
{
	insert_before(list, new_element);
}

ELEM* delete(ELEM *elem_to_delete)
{
	ELEM* prev = elem_to_delete->prev;
	ELEM* next = elem_to_delete->next;
	next->prev = prev;
	prev->next = next;
	elem_to_delete->prev = elem_to_delete->next = NULL;
	return elem_to_delete;
}

ELEM* delete_before(ELEM* before_this)
{
	return delete(before_this->prev);
}

void free_list(ELEM* list)
{
	ELEM *step = list->next;
	while(step != list)
	{
		step = step->next;
		free(step->prev);
	}
	free(list);
}

void wriie_to_list(ELEM *list, int len)
{
	for(int i = 0; i < len; i++)
	{
		ELEM *new_element = malloc(sizeof(ELEM));
		scanf("%d", &new_element->value);
		insert_after_tail(list, new_element);
	}
}

void insert_sort_list(ELEM* list)
{
	ELEM *inserted = list->next->next;
	while(inserted != list)
	{
		ELEM *next = inserted->next;
		ELEM *check = inserted->prev;
		while(check != list and check->value > inserted->value)
			check = check->prev;

		insert_after(check, delete_before(next));
		inserted = next;
	}
}

int main(int argc, char const *argv[])
{
	int len;
	scanf("%d", &len);
	ELEM *list = init_double_linked_list();
	wriie_to_list(list, len);
	// print_list(list, 1);
	insert_sort_list(list);
	print_list(list, 0);
	free_list(list);
	return 0;
}