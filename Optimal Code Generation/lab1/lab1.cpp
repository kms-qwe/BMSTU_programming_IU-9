#undef _FORTIFY_SOURCE

#include <cstdio>
#include <cstdint>

#include <gcc-plugin.h>
#include <coretypes.h>
#include <context.h>
#include <basic-block.h>
#include <tree-core.h>
#include <tree.h>
#include <tree-pass.h>
#include <gimple.h>
#include <gimple-ssa.h>
#include <gimple-iterator.h>
#include <gimple-expr.h>
#include <cfg.h>

std::int32_t plugin_is_GPL_compatible;

static const pass_data kPassData = {
    GIMPLE_PASS,      // type
    "bb-ir-dump",     // name
    OPTGROUP_NONE,    // optinfo_flags
    TV_NONE,          // tv_id
    PROP_gimple_any,  // properties_required
    0,                // properties_provided
    0,                // properties_destroyed
    0,                // todo_flags_start
    0                 // todo_flags_finish
};

class BbIrDumpPass final : public gimple_opt_pass {
public:
  explicit BbIrDumpPass(gcc::context *ctx) : gimple_opt_pass(kPassData, ctx) {}

  std::uint32_t execute(function *fn) override;
  BbIrDumpPass *clone() override { return this; }
};

static const char *safe_ident(tree t, const char *fallback) {
  if (!t || !DECL_NAME(t))
    return fallback;
  return IDENTIFIER_POINTER(DECL_NAME(t));
}

static void print_op_symbol(enum tree_code code) {
  switch (code) {
  case PLUS_EXPR:      printf("+");  break;
  case MINUS_EXPR:     printf("-");  break;
  case MULT_EXPR:      printf("*");  break;
  case RDIV_EXPR:      printf("/");  break;
  case BIT_IOR_EXPR:   printf("|");  break;
  case BIT_NOT_EXPR:   printf("!");  break;
  case TRUTH_AND_EXPR: printf("&&"); break;
  case TRUTH_OR_EXPR:  printf("||"); break;
  case TRUTH_NOT_EXPR: printf("!");  break;
  case LT_EXPR:        printf("<");  break;
  case LE_EXPR:        printf("<="); break;
  case GT_EXPR:        printf(">");  break;
  case GE_EXPR:        printf(">="); break;
  case EQ_EXPR:        printf("=="); break;
  case NE_EXPR:        printf("!="); break;
  default:
    printf("op(%d)", (int)code);
    break;
  }
}

static void print_tree_node(tree t);

static void print_phi_as_value(tree ssa) {
  gimple *def = SSA_NAME_DEF_STMT(ssa);

  printf("(%s__v%d = GIMPLE_PHI(",
         SSA_NAME_IDENTIFIER(ssa)
             ? IDENTIFIER_POINTER(SSA_NAME_IDENTIFIER(ssa))
             : "ssa_name",
         SSA_NAME_VERSION(ssa));

  for (unsigned i = 0; i < gimple_phi_num_args(def); ++i) {
    print_tree_node(gimple_phi_arg_def(def, i));
    if (i + 1 != gimple_phi_num_args(def))
      printf(", ");
  }

  printf("))");
}

static void print_tree_node(tree t) {
  if (!t) {
    printf("<null>");
    return;
  }

  switch (TREE_CODE(t)) {
  case INTEGER_CST:
    printf("%lld", (long long) TREE_INT_CST_LOW(t));
    break;

  case STRING_CST:
    printf("\"%s\"", TREE_STRING_POINTER(t));
    break;

  case LABEL_DECL:
    printf("%s:", DECL_NAME(t) ? IDENTIFIER_POINTER(DECL_NAME(t)) : "label_decl");
    break;

  case VAR_DECL:
    printf("%s", safe_ident(t, "var_decl"));
    break;

  case CONST_DECL:
    printf("%s", safe_ident(t, "const_decl"));
    break;

  case ARRAY_REF:
    print_tree_node(TREE_OPERAND(t, 0));
    printf("[");
    print_tree_node(TREE_OPERAND(t, 1));
    printf("]");
    break;

  case MEM_REF:
    printf("((typeof(");
    print_tree_node(TREE_OPERAND(t, 1));
    printf("))");
    print_tree_node(TREE_OPERAND(t, 0));
    printf(")");
    break;

  case SSA_NAME: {
    gimple *def = SSA_NAME_DEF_STMT(t);
    if (def && gimple_code(def) == GIMPLE_PHI) {
      print_phi_as_value(t);
    } else {
      printf("%s__v%d",
             SSA_NAME_IDENTIFIER(t)
                 ? IDENTIFIER_POINTER(SSA_NAME_IDENTIFIER(t))
                 : "ssa_name",
             SSA_NAME_VERSION(t));
    }
    break;
  }

  default:
    printf("tree_code(%d)", (int)TREE_CODE(t));
    break;
  }
}

static void dump_assign(gimple *stmt) {
  printf("      GIMPLE_ASSIGN(%d) { ", (int)GIMPLE_ASSIGN);

  const unsigned nops = gimple_num_ops(stmt);
  if (nops == 2) {
    print_tree_node(gimple_assign_lhs(stmt));
    printf(" = ");
    print_tree_node(gimple_assign_rhs1(stmt));
  } else if (nops == 3) {
    print_tree_node(gimple_assign_lhs(stmt));
    printf(" = ");
    print_tree_node(gimple_assign_rhs1(stmt));
    printf(" ");
    print_op_symbol(gimple_assign_rhs_code(stmt));
    printf(" ");
    print_tree_node(gimple_assign_rhs2(stmt));
  } else {
    printf("unsupported_assign_form");
  }

  printf(" }\n");
}

static void dump_call(gimple *stmt) {
  printf("      GIMPLE_CALL(%d) { ", (int)GIMPLE_CALL);

  tree lhs = gimple_call_lhs(stmt);
  if (lhs) {
    print_tree_node(lhs);
    printf(" = ");
  }

  tree callee = gimple_call_fndecl(stmt);
  printf("%s(", callee ? fndecl_name(callee) : "<indirect-call>");

  for (unsigned i = 0; i < gimple_call_num_args(stmt); ++i) {
    print_tree_node(gimple_call_arg(stmt, i));
    if (i + 1 != gimple_call_num_args(stmt))
      printf(", ");
  }

  printf(") }\n");
}

static void dump_cond(gimple *stmt) {
  printf("      GIMPLE_COND(%d) { ", (int)GIMPLE_COND);
  print_tree_node(gimple_cond_lhs(stmt));
  printf(" ");
  print_op_symbol(gimple_cond_code(stmt));
  printf(" ");
  print_tree_node(gimple_cond_rhs(stmt));
  printf(" }\n");
}

static void dump_label([[maybe_unused]] gimple *stmt) {
  printf("      GIMPLE_LABEL(%d)\n", (int)GIMPLE_LABEL);
}

static void dump_ret([[maybe_unused]] gimple *stmt) {
  printf("      GIMPLE_RETURN(%d)\n", (int)GIMPLE_RETURN);
}

static void dump_phi_nodes(basic_block bb) {
  for (gphi_iterator pit = gsi_start_phis(bb); !gsi_end_p(pit); gsi_next(&pit)) {
    gphi *phi = pit.phi();
    printf("      GIMPLE_PHI(%d) { ", (int)GIMPLE_PHI);

    print_tree_node(gimple_phi_result(phi));
    printf(" = PHI(");

    for (unsigned i = 0; i < gimple_phi_num_args(phi); ++i) {
      print_tree_node(gimple_phi_arg_def(phi, i));
      printf(" from bb%d", gimple_phi_arg_edge(phi, i)->src->index);
      if (i + 1 != gimple_phi_num_args(phi))
        printf(", ");
    }

    printf(") }\n");
  }
}

static void dump_edges(vec<edge, va_gc> *edges, bool is_pred) {
  edge e;
  edge_iterator ei;
  FOR_EACH_EDGE(e, ei, edges) {
    printf("%d ", is_pred ? e->src->index : e->dest->index);
  }
}

std::uint32_t BbIrDumpPass::execute(function *fn) {
  printf("function: \"%s\"\n", function_name(fn));

  basic_block bb;
  FOR_ALL_BB_FN(bb, fn) {
    printf("  basic block: %d\n", bb->index);

    printf("    preds: ");
    dump_edges(bb->preds, true);
    printf("\n");

    printf("    succs: ");
    dump_edges(bb->succs, false);
    printf("\n");

    printf("    statements:\n");

    dump_phi_nodes(bb);

    for (gimple_stmt_iterator gsi = gsi_start_bb(bb); !gsi_end_p(gsi); gsi_next(&gsi)) {
      gimple *stmt = gsi_stmt(gsi);

      switch (gimple_code(stmt)) {
      case GIMPLE_ASSIGN:
        dump_assign(stmt);
        break;
      case GIMPLE_CALL:
        dump_call(stmt);
        break;
      case GIMPLE_COND:
        dump_cond(stmt);
        break;
      case GIMPLE_LABEL:
        dump_label(stmt);
        break;
      case GIMPLE_RETURN:
        dump_ret(stmt);
        break;
      default:
        break;
      }
    }
  }

  return 0;
}

static register_pass_info kPassInfo = {
    new BbIrDumpPass(g),
    "ssa",
    1,
    PASS_POS_INSERT_AFTER
};

static plugin_info kPluginInfo = {
    "0.2.0",
    "manual GIMPLE/CFG printer"
};

int plugin_init(struct plugin_name_args *info,
                [[maybe_unused]] struct plugin_gcc_version *ver) {
  register_callback(info->base_name, PLUGIN_INFO, nullptr, &kPluginInfo);
  register_callback(info->base_name, PLUGIN_PASS_MANAGER_SETUP, nullptr, &kPassInfo);
  return 0;
}
