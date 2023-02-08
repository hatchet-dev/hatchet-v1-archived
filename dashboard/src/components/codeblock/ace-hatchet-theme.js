import ace from "brace";

ace["define"](
  "ace/theme/hatchet",
  ["require", "exports", "module", "ace/lib/dom"],
  (acequire, exports) => {
    exports.isDark = true;
    exports.cssClass = "ace-hatchet";
    exports.cssText = `.ace-hatchet, div.ace_content, div.ace_line, div.ace_gutter-cell {\
    font-family: monospace;
    font-size: 14px;
    }
    
    .ace-hatchet {
    background-color: #141321;
    }
    .ace_identifier {
        color: #7AA6DA;
    }

    .ace_paren, .ace_punctuation {
        color: #D272FF;
    }

    .ace-hatchet .ace_gutter {
    background: #04061c;
    color: #929292;
    }

    .ace_line {
        color: #D272FF;
    }

    .ace-hatchet .ace_print-margin {
    width: 1px;
    background: #04061c;
    }
    .ace-hatchet .ace_cursor {
    color: #bfc7d5;
    }
    .ace-hatchet .ace_marker-layer .ace_selection {
    background: #04061c;
    }
    .ace-hatchet.ace_multiselect .ace_selection.ace_start {
    box-shadow: 0 0 3px 0px #191919;
    }
    .ace-hatchet .ace_marker-layer .ace_step {
    background: rgb(102, 82, 0);
    }
    .ace-hatchet .ace_marker-layer .ace_bracket {
    margin: -1px 0 0 -1px;
    border: 1px solid #BFBFBF;
    }
    .ace-hatchet .ace_marker-layer .ace_active-line {
    background: rgba(215, 215, 215, 0.031);
    }
    .ace-hatchet .ace_gutter-active-line {
    background-color: rgba(215, 215, 215, 0.031);
    }
    .ace-hatchet .ace_marker-layer .ace_selected-word {
    border: 1px solid #424242;
    }
    .ace-hatchet .ace_invisible {
    color: #343434;
    }
    .ace-hatchet .ace_keyword,
    .ace-hatchet .ace_meta,
    .ace-hatchet .ace_storage,
    .ace-hatchet .ace_storage.ace_type,
    .ace-hatchet .ace_support.ace_type {
    color: #D272FF;
    }
    .ace-hatchet .ace_keyword.ace_operator {
    color: #D272FF;
    }
    .ace-hatchet .ace_constant.ace_character,
    .ace-hatchet .ace_constant.ace_language,
    .ace-hatchet .ace_constant.ace_numeric,
    .ace-hatchet .ace_keyword.ace_other.ace_unit,
    .ace-hatchet .ace_support.ace_constant,
    .ace-hatchet .ace_variable.ace_parameter {
    color: #D272FF;
    }
    .ace-hatchet .ace_constant.ace_other {
    color: gold
    }
    .ace-hatchet .ace_invalid {
    color: yellow;
    background-color: red
    }
    .ace-hatchet .ace_invalid.ace_deprecated {
    color: #CED2CF;
    background-color: #B798BF
    }
    .ace-hatchet .ace_fold {
    background-color: #7AA6DA;
    border-color: #DEDEDE
    }
    .ace-hatchet .ace_entity.ace_name.ace_function,
    .ace-hatchet .ace_support.ace_function,
    .ace-hatchet .ace_variable {
    color: #7AA6DA
    }
    .ace-hatchet .ace_support.ace_class,
    .ace-hatchet .ace_support.ace_type {
    color: #E7C547
    }
    .ace-hatchet .ace_heading,
    .ace-hatchet .ace_string {
    color: #23F0C3;
    }
    .ace-hatchet .ace_entity.ace_name.ace_tag,
    .ace-hatchet .ace_entity.ace_other.ace_attribute-name,
    .ace-hatchet .ace_meta.ace_tag,
    .ace-hatchet .ace_string.ace_regexp,
    .ace-hatchet .ace_variable {
    color: #23F0C3;
    }
    .ace-hatchet .ace_comment {
    color: #717AA8;
    }
    .ace-hatchet .ace_indent-guide {
    background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAACCAYAAACZgbYnAAAAEklEQVQImWNgYGBgYLBWV/8PAAK4AYnhiq+xAAAAAElFTkSuQmCC) right repeat-y;
    }`;

    var dom = acequire("../lib/dom");
    dom.importCssString(exports.cssText, exports.cssClass);
  }
);
