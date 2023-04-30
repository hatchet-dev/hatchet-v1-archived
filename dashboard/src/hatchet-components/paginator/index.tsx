import { HoverableMaterialIcon } from "../globals";
import React from "react";
import { PaginatorText, PaginatorWrapper } from "./styles";

type Props = {
  curr_page: number;
  last_page: number;
  cursor_forward: () => void;
  cursor_backward: () => void;
};

const Paginator: React.FC<Props> = ({
  curr_page,
  last_page,
  cursor_forward,
  cursor_backward,
}) => {
  return (
    <PaginatorWrapper>
      <HoverableMaterialIcon
        className="material-icons"
        onClick={() => cursor_backward()}
      >
        keyboard_arrow_left
      </HoverableMaterialIcon>
      <PaginatorText>
        {curr_page} of {last_page}
      </PaginatorText>
      <HoverableMaterialIcon
        className="material-icons"
        onClick={() => cursor_forward()}
      >
        keyboard_arrow_right
      </HoverableMaterialIcon>
    </PaginatorWrapper>
  );
};

export default Paginator;
