import { MaterialIcon, Relative, Span } from "components/globals";
import React, { useEffect, useRef, useState } from "react";
import {
  Dropdown,
  DropdownWrapper,
  ScrollableWrapper,
  SelectorPlaceholder,
  StyledSelection,
  StyledSelector,
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
  select?: (option: Selection) => void;
};

const Selector: React.FC<Props> = ({
  placeholder,
  placeholder_icon,
  placeholder_material_icon,
  options,
  select,
}) => {
  const [selection, setSelection] = useState<Selection>();

  const [expanded, setExpanded] = useState(false);

  const wrapperRef = useRef<HTMLInputElement>(null);
  const parentRef = useRef<HTMLInputElement>(null);

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
        <DropdownWrapper>
          <Dropdown ref={wrapperRef}>
            {options.length > 0 ? (
              <ScrollableWrapper>
                {options.map((option) => {
                  return (
                    <StyledSelection onClick={() => onClickSelection(option)}>
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
    if (selection) {
      return (
        <SelectorPlaceholder>
          {selection.icon ? (
            <img src={selection.icon} />
          ) : (
            <MaterialIcon className="material-icons">
              {selection.material_icon}
            </MaterialIcon>
          )}
          <div>{selection.label}</div>
          <i className="material-icons">expand_more</i>
        </SelectorPlaceholder>
      );
    }

    return (
      <SelectorPlaceholder>
        {placeholder_icon ? (
          <img src={placeholder_icon} />
        ) : (
          <MaterialIcon className="material-icons">
            {placeholder_material_icon}
          </MaterialIcon>
        )}

        <div>{placeholder}</div>
        <i className="material-icons">expand_more</i>
      </SelectorPlaceholder>
    );
  };

  return (
    <Relative>
      <StyledSelector onClick={() => setExpanded(!expanded)} ref={parentRef}>
        {renderPlaceholder()}
      </StyledSelector>
      {renderDropdown()}
    </Relative>
  );
};

export default Selector;
