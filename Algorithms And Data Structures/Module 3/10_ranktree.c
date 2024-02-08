#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <stdbool.h>
#include <string.h>
#include <time.h>
#include <limits.h>
#define COMMNAD_SIZE 10
#define STR_SIZE 11

typedef struct TOP
{
	int key;
	char *str;
	struct TOP *parent, *left, *right;
	int count;
} TOP;

TOP* init_binary_search_tree()
{
	TOP *tree = NULL;
	return tree;
}

bool map_empty(TOP *tree)
{
	return tree == NULL;
}

TOP* descend(TOP *tree, int searched_key)
{
	TOP *step = tree;
	while(step != NULL and step->key != searched_key)
		step = searched_key < step->key ? step->left : step->right;
	return step;
}

bool map_search(TOP *tree, int searched_key)
{
	return descend(tree, searched_key) != NULL;
}

char* lookup(TOP *tree, int searched_key)
{
	TOP *step = descend(tree, searched_key);
	if(step == NULL)
		printf("ERROR: NO SUCH KEY (lookup)\n");
	return step->str;
}

void increase_count(TOP *new_top)
{
	TOP *step = new_top->parent;
	while(step != NULL)
	{
		step->count++;
		step = step->parent;
	}
}

void insert(TOP **tree, int new_key, char *new_str)
{
	TOP *new_top = malloc(sizeof(TOP));
	new_top->key = new_key;
	new_top->str = malloc(STR_SIZE);
	strcpy(new_top->str, new_str);
	new_top->parent = new_top->left = new_top->right = NULL;
	new_top->count = 1;
	if(map_empty(*tree))
		*tree = new_top;
	else
	{
		TOP *step = *tree;
		while(true)
		{
			if(step->key == new_key)
				printf("ERROR: KEY IS TAKEN (insert)\n");
			if(new_key < step->key)
			{
				if(step->left == NULL)
				{
					step->left = new_top;
					new_top->parent = step;
					break;
				}
				step = step->left;
			}
			else
			{
				if(step->right == NULL)
				{
					step->right = new_top;
					new_top->parent = step;
					break;
				}
				step = step->right;
			}
		}
	}
	increase_count(new_top);
}
void print_if_right(TOP *top)
{
	if(top->right == NULL)
		printf("right = NULL, ");
	else
		printf("right_key = %d, ", top->right->key);
}
void print_if_left(TOP *top)
{
	if(top->left == NULL)
		printf("left = NULL, ");
	else
		printf("left_key = %d, ", top->left->key);
}
void print_if_parent(TOP *top)
{
	if(top->parent == NULL)
		printf("parent = NULL, ");
	else
		printf("parent_key = %d, ", top->parent->key);
}
void print_tree(TOP *tree)
{
	if(map_empty(tree))
		printf("TREE iS EMPTY\n");
	else
	{
		printf("count = %d, key = %d, ", tree->count, tree->key);
		print_if_parent(tree);
		print_if_left(tree);
		print_if_right(tree);
		printf("\nstr = |%s|\n", tree->str);
		print_tree(tree->left);
		print_tree(tree->right);
	}

}
void print_top(TOP *tree)
{
	if(map_empty(tree))
		printf("TREE iS EMPTY\n");
	else
	{
		printf("count = %d, key = %d, ", tree->count, tree->key);
		print_if_parent(tree);
		print_if_left(tree);
		print_if_right(tree);
		printf("\nstr = |%s|\n", tree->str);
	}

}
void replace_node(TOP **tree, TOP *old_top, TOP *new_top)
{
	if(old_top == *tree)
	{
		*tree = new_top;
		if(new_top != NULL) new_top->parent = NULL;
	}
	else
	{
		TOP *parent_old_top = old_top->parent;
		if(new_top != NULL) new_top->parent = parent_old_top;
		if(parent_old_top->left == old_top) 
			parent_old_top->left = new_top;
		else 
			parent_old_top->right = new_top;
	}
}

TOP* minimum(TOP *tree)
{
	if(map_empty(tree))
		return NULL;
	TOP *step = tree;
	while(step->left != NULL)
		step = step->left;
	return step;
}

TOP* succ(TOP *prev)
{
	if(prev->right != NULL)
		return minimum(prev->right);
	TOP *next = prev->parent;
	while(next != NULL and prev == next->right)
	{
		prev = next;
		next = next->parent;
	}
	return next;
}

void decrease_count(TOP *deleted_top)
{
	TOP *step = deleted_top->parent;
	while(step != NULL)
	{
		step->count--;
		step = step->parent;
	}
}

void free_top(TOP *top)
{
	free(top->str);
	free(top);
}



void delete(TOP **tree, int delete_key)
{
	TOP *step = descend(*tree, delete_key);
	if(step == NULL) printf("ERROR: NO SUCH KEY (delete)\n");
	if(step->left == NULL and step->right == NULL) 
	{
		decrease_count(step);
		replace_node(tree, step, NULL);
	}
	else if(step->left == NULL)
	{
		decrease_count(step);
		replace_node(tree, step, step->right);
	}
	else if(step->right == NULL)
	{
		decrease_count(step);
		replace_node(tree, step, step->left);
	}
	else
	{
		TOP *next = succ(step);
		decrease_count(next);
		next->count = step->count;
		replace_node(tree, next, next->right);
		step->left->parent = next;
		next->left = step->left;
		if(step->right != NULL) step->right->parent = next;
		next->right = step->right;
		replace_node(tree, step, next);
	}
	free_top(step);
}

char* search_by_index(TOP* tree, int searched_index, int low, int high)
{
	if(map_empty(tree)) printf("ERROR: tree is empty (search_by_index)\n");
	if(tree->right != NULL and searched_index >= high + 1 - tree->right->count)
		return search_by_index(tree->right, searched_index, high + 1 - tree->right->count, high);
	if(tree->left != NULL and searched_index <= low - 1 + tree->left->count)
		return search_by_index(tree->left, searched_index, low, low - 1 + tree->left->count);
	return tree->str;
}

void scan_command(char* command)
{
	fgets(command, COMMNAD_SIZE, stdin);
	command[strlen(command) - 1] = '\0'; 
}

void free_tree(TOP **tree)
{
	while(*tree != NULL)
		delete(tree, (*tree)->key);
		
}

void run()
{
	char* command = calloc(COMMNAD_SIZE, sizeof(char));
	char* str = calloc(STR_SIZE, sizeof(char));
	int key, index;
	TOP *tree = init_binary_search_tree();
	scanf("%s", command);
	while(strcmp(command, "END"))
	{
		if(!strcmp(command, "INSERT"))
		{
			scanf("%d", &key);
			scanf("%s", str);
			insert(&tree, key, str);
		}
		if(!strcmp(command, "LOOKUP"))
		{
			scanf("%d", &key);
			printf("%s\n", lookup(tree, key));
		}
		if(!strcmp(command, "DELETE"))
		{
			scanf("%d", &key);
			delete(&tree, key);
		}
		if(!strcmp(command, "SEARCH"))
		{
			scanf("%d", &index);
			printf("%s\n", search_by_index(tree, index + 1, 1, tree->count));
		}
		// printf("====================NEW ITERATION====================\n");
		// print_tree(tree);
		// printf("====================     END     ====================\n");
		scanf("%s", command);
	}

	free(command);
	free(str);
	free_tree(&tree);
}

int main(int argc, char const *argv[])
{
	run();
	return 0;
}