import { MaterialIcon } from "components/globals";
import React from "react";
import { FilterText, FilterWrapper } from "./styles";

export type Props = {};

const Filter: React.FC<Props> = () => {
  return (
    <FilterWrapper>
      <MaterialIcon className="material-icons">filter_list</MaterialIcon>
      <FilterText>Filter</FilterText>
    </FilterWrapper>
  );
};

export default Filter;
