import {
  MaterialIcon,
  SmallSpan,
  Span,
  Placeholder,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useRef, useState } from "react";
import usePrevious from "../../shared/hooks/useprevious";
import {
  Dropdown,
  DropdownWrapper,
  InnerMultiSelectorPlaceholder,
  ScrollableWrapper,
  MultiSelectorPlaceholder,
  StyledSelection,
  StyledMultiSelector,
  StyledMultiSelectorWrapper,
} from "./styles";

export type Selection = {
  label: string;
  value: string;
  icon?: string;
  material_icon?: string;
  checked?: boolean;
};

type Props = {
  placeholder: string;
  placeholder_icon?: string;
  placeholder_material_icon?: string;
  options: Selection[];
  orientation?: "horizontal" | "vertical";
  icon_style?: "circular" | "square";
  option_alignment?: "left" | "right";
  fill_selection?: boolean;
  select?: (option: Selection) => void;
  reset?: number;
  selector_wrapper_overrides?: string;
  selector_overrides?: string;
};

const MultiSelector: React.FC<Props> = ({
  placeholder,
  placeholder_icon,
  placeholder_material_icon,
  options,
  orientation = "horizontal",
  option_alignment = "left",
  icon_style = "circular",
  fill_selection = true,
  select,
  reset,
  selector_overrides,
  selector_wrapper_overrides,
}) => {
  const [selections, setSelections] = useState<Selection[]>();

  const [expanded, setExpanded] = useState(false);

  const wrapperRef = useRef<HTMLInputElement>(null);
  const parentRef = useRef<HTMLInputElement>(null);
  const prevReset = usePrevious(reset);

  useEffect(() => {
    if (reset != prevReset) {
      setSelections(null);
    }
  }, [reset]);

  useEffect(() => {
    document.addEventListener("mousedown", handleClickOutside.bind(this));
    return () =>
      document.removeEventListener("mousedown", handleClickOutside.bind(this));
  }, []);

  const handleClickOutside = (event: any) => {
    if (
      wrapperRef &&
      wrapperRef.current &&
      !wrapperRef.current.contains(event.target) &&
      parentRef &&
      parentRef.current &&
      !parentRef.current.contains(event.target)
    ) {
      setExpanded(false);
    }
  };

  const onClickSelection = (selection: Selection) => {
    // setSelections([...selections, selection]);
    // setExpanded(false);
    select && select(selection);
  };

  const renderDropdown = () => {
    if (expanded) {
      return (
        <DropdownWrapper align={option_alignment} orientation={orientation}>
          <Dropdown ref={wrapperRef}>
            {options.length > 0 ? (
              <ScrollableWrapper>
                {options.map((option) => {
                  return (
                    <StyledSelection
                      onClick={() => onClickSelection(option)}
                      key={option.value}
                      icon_style={icon_style}
                    >
                      {option.icon ? (
                        <img src={option.icon} />
                      ) : (
                        <MaterialIcon className="material-icons">
                          {option.checked
                            ? "check_box"
                            : "check_box_outline_blank"}
                        </MaterialIcon>
                      )}
                      <div>{option.label}</div>
                    </StyledSelection>
                  );
                })}
              </ScrollableWrapper>
            ) : (
              <Placeholder>
                <SmallSpan>No options found</SmallSpan>
              </Placeholder>
            )}
          </Dropdown>
        </DropdownWrapper>
      );
    }
  };

  const renderPlaceholder = () => {
    // if (fill_selection && selection) {
    //   return (
    //     <MultiSelectorPlaceholder>
    //       <InnerMultiSelectorPlaceholder icon_style={icon_style}>
    //         {selection.icon ? (
    //           <img src={selection.icon} />
    //         ) : (
    //           <MaterialIcon className="material-icons">
    //             {selection.material_icon}
    //           </MaterialIcon>
    //         )}
    //         <div>{selection.label}</div>
    //       </InnerMultiSelectorPlaceholder>

    //       <i className="material-icons">
    //         {orientation == "horizontal" ? "expand_more" : "chevron_right"}
    //       </i>
    //     </MultiSelectorPlaceholder>
    //   );
    // }
    const numSelected = options.filter((option) => option.checked).length;

    if (numSelected > 0) {
      return (
        <MultiSelectorPlaceholder>
          <InnerMultiSelectorPlaceholder icon_style={icon_style}>
            {placeholder_icon ? (
              <img src={placeholder_icon} />
            ) : (
              <MaterialIcon className="material-icons">
                {placeholder_material_icon}
              </MaterialIcon>
            )}

            <div>{numSelected} Modules</div>
          </InnerMultiSelectorPlaceholder>
          <i className="material-icons">
            {orientation == "horizontal" ? "expand_more" : "chevron_right"}
          </i>
        </MultiSelectorPlaceholder>
      );
    }

    return (
      <MultiSelectorPlaceholder>
        <InnerMultiSelectorPlaceholder icon_style={icon_style}>
          {placeholder_icon ? (
            <img src={placeholder_icon} />
          ) : (
            <MaterialIcon className="material-icons">
              {placeholder_material_icon}
            </MaterialIcon>
          )}

          <div>{placeholder}</div>
        </InnerMultiSelectorPlaceholder>
        <i className="material-icons">
          {orientation == "horizontal" ? "expand_more" : "chevron_right"}
        </i>
      </MultiSelectorPlaceholder>
    );
  };

  return (
    <StyledMultiSelectorWrapper
      orientation={orientation}
      overrides={selector_wrapper_overrides}
    >
      <StyledMultiSelector
        onClick={() => {
          setExpanded(!expanded);
        }}
        ref={parentRef}
        orientation={orientation}
        overrides={selector_overrides}
      >
        {renderPlaceholder()}
      </StyledMultiSelector>
      {renderDropdown()}
    </StyledMultiSelectorWrapper>
  );
};

export default MultiSelector;
