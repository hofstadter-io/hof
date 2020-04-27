package schema

#TemplateConfigReplacible: {
  // System params
  TemplateSystem: *"." | "golang" | "raymond"

  //
  // Template delimiters
  //
  //   these are for advanced usage, you shouldn't have to modify them normally

  // Alt and Swap Delims,
  //   becuase the defaulttemplate systems use `{{` and `}}`
  //   and you may choose to use other delimiters, but the lookup system is still based on the template system
  //   and if you want to preserve those, we need three sets of delimiters
  AltDelims:  bool | *false
  SwapDelims: bool | *false

  // The default delimiters
  // You should change these when using alternative style like jinjas {% ... %}
  // They also need to be different when using the swap system
  LHS_D: LHS2_D
  RHS_D: RHS2_D
  LHS2_D: string | *"."
  RHS2_D: string | *"."
  LHS3_D: string | *"."
  RHS3_D: string | *"."

  // These are the same as the default becuase
  // the current template systems require these.
  //   So these should really never change or be overriden until there is a new template system
  //     supporting setting the delimiters dynamicalldelimiters dynamicallyy
  LHS_S: LHS2_S
  RHS_S: RHS2_S
  LHS2_S: string | *"."
  RHS2_S: string | *"."
  LHS3_S: string | *"."
  RHS3_S: string | *"."

  // The temporary delims to replace swap with while also swapping
  // the defaults you set to the swap that is required by the current templet systems
  // You need this when you are double templating a file and the top-level system is not the default
  LHS_T: LHS2_T
  RHS_T: RHS2_T
  LHS2_T: string | *"."
  RHS2_T: string | *"."
  LHS3_T: string | *"."
  RHS3_T: string | *"."
}

#DefaultTemplateConfig: {
  // Include Common attributes
  // System params
  TemplateSystem: *"golang" | "raymond"

  //
  // Template delimiters
  //
  //   these are for advanced usage, you shouldn't have to modify them normally

  // Alt and Swap Delims,
  //   becuase the defaulttemplate systems use `{{` and `}}`
  //   and you may choose to use other delimiters, but the lookup system is still based on the template system
  //   and if you want to preserve those, we need three sets of delimiters
  AltDelims:  bool | *false
  SwapDelims: bool | *false

  // The default delimiters
  // You should change these when using alternative style like jinjas {% ... %}
  // They also need to be different when using the swap system
  LHS_D: LHS2_D
  RHS_D: RHS2_D
  LHS2_D: string | *"{{"
  RHS2_D: string | *"}}"
  LHS3_D: string | *"{{{"
  RHS3_D: string | *"}}}"

  // These are the same as the default becuase
  // the current template systems require these.
  //   So these should really never change or be overriden until there is a new template system
  //     supporting setting the delimiters dynamicalldelimiters dynamicallyy
  LHS_S: LHS2_S
  RHS_S: RHS2_S
  LHS2_S: string | *"{{"
  RHS2_S: string | *"}}"
  LHS3_S: string | *"{{{"
  RHS3_S: string | *"}}}"

  // The temporary delims to replace swap with while also swapping
  // the defaults you set to the swap that is required by the current templet systems
  // You need this when you are double templating a file and the top-level system is not the default
  LHS_T: LHS2_T
  RHS_T: RHS2_T
  LHS2_T: string | *"#_hof_l2_#"
  RHS2_T: string | *"#_hof_r2_#"
  LHS3_T: string | *"#_hof_l3_#"
  RHS3_T: string | *"#_hof_r3_#"
}

