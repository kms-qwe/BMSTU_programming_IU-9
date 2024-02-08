#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <stdbool.h>
#include <string.h>
#include <time.h>
#include <limits.h>
#define STR_SIZE 1100

typedef struct ELEM
{
	struct ELEM *next;
	char *word;
} ELEM;

ELEM* init_single_linked_list()
{
	ELEM* list = NULL;
	return list;
}

bool list_empty(ELEM *list)
{
	return list == NULL;
}

void print_list(ELEM *list, int key)
{
	if(list_empty(list)) printf("ERROR: list empty (print_list)\n");
	if(key)
	{
		int i = 0;
		ELEM *step = list;
		while(step)
		{
			printf("list[%d] = |%s|\n", i, step->word);
			step = step->next;
			i++;
		}
	}
	else
	{
		ELEM *step = list;
		while(step)
		{
		if(step->next) printf("%s ", step->word);
		else printf("%s\n", step->word);
		step = step->next;
		}
	}
}

int list_length(ELEM *list)
{
	int len = 0;
	ELEM *step = list;
	while(step)
	{
		len++;
		step = step->next;
	}
	return len;
}

void insert_after(ELEM *after_this, ELEM *insert_this)
{
	if(!after_this) printf("ERROR: after_this is NULL (insert_after)\n");
	ELEM *after_after = after_this->next;
	after_this->next = insert_this;
	insert_this->next = after_after;
}

void insert_before_head(ELEM *list, ELEM *new_element)
{
	if(list_empty(list)) printf("ERROR: list empty (insert_before_head)\n");
	new_element->next = list;
	list = new_element;
}

ELEM* delete_after(ELEM *delete_after_this)
{
	if(!delete_after_this) printf("ERROR: delete_after_this is NULL (delete_after)\n");
	ELEM *elem_to_delete = delete_after_this->next;
	delete_after_this = elem_to_delete->next;
	elem_to_delete->next = NULL;
	return elem_to_delete;
}

ELEM* delete_head(ELEM* list)
{
	if(list_empty(list)) printf("ERROR: list empty (delete_head)\n");
	ELEM *head = list;
	list = head->next;
	head->next = NULL;
	return head;
}

int wcount(char *s)
{
    int count = 0, len = strlen(s);
    if (len == 0)
        return 0;
    if (1 - (s[0] == ' ' or s[0] == '\n' or s[0] == '\t'))
        count += 1;
    for (int i = 1; i < len; i++)
    {
        if ((s[i-1] == ' ' or s[i-1] == '\n' or s[i-1] == '\t') 
        and (1 - (s[i] == ' ' or s[i] == '\n' or s[i] == '\t')))
            count += 1;
    }
    return count;
}

char** split(char *source, int len_split)
{
	char **split_array = malloc(len_split * sizeof(char *));
	for(int i = 0; i < len_split; i++)
		split_array[i] = calloc(STR_SIZE/len_split + 10, sizeof(char));

	int source_index = 0, split_index = 0, len_source = strlen(source);
	while(split_index < len_split and source_index < len_source)
	{
		int word_index = 0;
		while(source_index < len_source and 
			source[source_index] != '\t' and source[source_index] != '\n'
			and source[source_index] != ' ')
		{
			split_array[split_index][word_index] = source[source_index];
			word_index++;
			source_index++;
		}
		if(word_index != 0) split_index++;
		source_index++;
	}
	return split_array;
}

void free_array(char **s, int len)
{
	for(int i = 0; i < len; i++)
		free(s[i]);
	free(s);
}

void free_list(ELEM *list)
{
	ELEM *current = list;
	ELEM *next = list->next;
	while(next)
	{
		free(current->word);
		free(current);
		current = next;
		next = next->next;
	}
	free(current->word);
	free(current);
}

ELEM* write_to_list(char *source, int len_split)
{
	if(len_split == 0) printf("SOURCE IS EMPTY\n");
	char **split_array = split(source, len_split);

	int split_index = 0;
	ELEM *list = malloc(sizeof(ELEM));
	list->word = split_array[split_index++];
	list->next = NULL;
	ELEM *prev = list;
	while(split_index < len_split)
	{
		ELEM *new_element = malloc(sizeof(ELEM));
		new_element->word = split_array[split_index++];
		new_element->next = NULL;
		insert_after(prev, new_element);
		prev = new_element;
	}
	free(source);
	free(split_array);
	return list;
}

void swap(ELEM *a, ELEM *b)
{
	char *tmp = a->word;
	a->word = b->word;
	b->word = tmp;
}

int compare(char *a, char *b)
{
	int len_a = strlen(a), len_b = strlen(b);
	if(len_a > len_b) return 1;
	if(len_a < len_b) return -1;
	return 0; 
}

ELEM *get_i(ELEM *list, int i)
{
	if(list_empty(list)) printf("LIST IS EMPTY (get_i)\n");
	int k = 0;
	ELEM *list_i = list;
	while(k < i and list_i->next != NULL)
	{
		list_i = list_i->next;
		k++;
	}
	return list_i;
}

ELEM* bsort(ELEM *list)
{
	int len = list_length(list);
	// if(list_empty(list)) return list;
	int t = len - 1;
	while(t > 0)
	{
		int bound = t;
		t = 0;
		int i = 0;
		while(i < bound)
		{
			ELEM *list_i = get_i(list, i);
			if(compare(list_i->next->word, list_i->word) < 0)
			{
				swap(list_i->next, list_i);
				t = i;
			}
			i++;
		}
	}
	return list;
}

int main(int argc, char const *argv[])
{
	char *source = calloc(STR_SIZE, sizeof(char));
	fgets(source, STR_SIZE, stdin);
	int len_split = wcount(source);
	if(len_split == 0)
	{
		printf("\n");
		free(source);
		return 0;
	}
	ELEM *list = write_to_list(source, len_split);
	list = bsort(list);
	print_list(list, 0);
	free_list(list);
	return 0;
}