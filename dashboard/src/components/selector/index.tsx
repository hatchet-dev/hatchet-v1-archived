import { MaterialIcon, Span } from "components/globals";
import React, { useEffect, useRef, useState } from "react";
import usePrevious from "shared/hooks/useprevious";
import {
  Dropdown,
  DropdownWrapper,
  InnerSelectorPlaceholder,
  ScrollableWrapper,
  SelectorPlaceholder,
  StyledSelection,
  StyledSelector,
  StyledSelectorWrapper,
} from "./styles";

export type Selection = {
  label: string;
  value: string;
  icon?: string;
  material_icon?: string;
};

export type Props = {
  placeholder: string;
  placeholder_icon?: string;
  placeholder_material_icon?: string;
  options: Selection[];
  orientation?: "horizontal" | "vertical";
  option_alignment?: "left" | "right";
  fill_selection?: boolean;
  select?: (option: Selection) => void;
  reset?: number;
};

const Selector: React.FC<Props> = ({
  placeholder,
  placeholder_icon,
  placeholder_material_icon,
  options,
  orientation = "horizontal",
  option_alignment = "left",
  fill_selection = true,
  select,
  reset,
}) => {
  const [selection, setSelection] = useState<Selection>();

  const [expanded, setExpanded] = useState(false);

  const wrapperRef = useRef<HTMLInputElement>(null);
  const parentRef = useRef<HTMLInputElement>(null);
  const prevReset = usePrevious(reset);

  useEffect(() => {
    if (reset != prevReset) {
      setSelection(null);
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
    setSelection(selection);
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
                    >
                      {option.icon ? (
                        <img src={option.icon} />
                      ) : (
                        <MaterialIcon className="material-icons">
                          {option.material_icon}
                        </MaterialIcon>
                      )}
                      <div>{option.label}</div>
                    </StyledSelection>
                  );
                })}
              </ScrollableWrapper>
            ) : (
              <Span>No options found</Span>
            )}
          </Dropdown>
        </DropdownWrapper>
      );
    }
  };

  const renderPlaceholder = () => {
    if (fill_selection && selection) {
      return (
        <SelectorPlaceholder>
          <InnerSelectorPlaceholder>
            {selection.icon ? (
              <img src={selection.icon} />
            ) : (
              <MaterialIcon className="material-icons">
                {selection.material_icon}
              </MaterialIcon>
            )}
            <div>{selection.label}</div>
          </InnerSelectorPlaceholder>

          <i className="material-icons">
            {orientation == "horizontal" ? "expand_more" : "chevron_right"}
          </i>
        </SelectorPlaceholder>
      );
    }

    return (
      <SelectorPlaceholder>
        <InnerSelectorPlaceholder>
          {placeholder_icon ? (
            <img src={placeholder_icon} />
          ) : (
            <MaterialIcon className="material-icons">
              {placeholder_material_icon}
            </MaterialIcon>
          )}

          <div>{placeholder}</div>
        </InnerSelectorPlaceholder>
        <i className="material-icons">
          {orientation == "horizontal" ? "expand_more" : "chevron_right"}
        </i>
      </SelectorPlaceholder>
    );
  };

  return (
    <StyledSelectorWrapper orientation={orientation}>
      <StyledSelector
        onClick={() => {
          setExpanded(!expanded);
        }}
        ref={parentRef}
        orientation={orientation}
      >
        {renderPlaceholder()}
      </StyledSelector>
      {renderDropdown()}
    </StyledSelectorWrapper>
  );
};

export default Selector;
