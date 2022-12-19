import { HoverableMaterialIcon } from "components/globals";
import React from "react";
import { PaginatorText, PaginatorWrapper } from "./styles";

export type Props = {};

const Paginator: React.FC<Props> = () => {
  return (
    <PaginatorWrapper>
      <HoverableMaterialIcon className="material-icons">
        keyboard_arrow_left
      </HoverableMaterialIcon>
      <PaginatorText>1 of 25</PaginatorText>
      <HoverableMaterialIcon className="material-icons">
        keyboard_arrow_right
      </HoverableMaterialIcon>
    </PaginatorWrapper>
  );
};

export default Paginator;
