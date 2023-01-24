import React, { useState } from "react";
import { PaginationResponse } from "shared/api/generated/data-contracts";

function usePagination() {
  const [currentPage, setCurrentPage] = useState(1);
  const [maxPage, setMaxPage] = useState(1);

  function cursor_forward() {
    setCurrentPage((currentPage) => Math.min(currentPage + 1, maxPage));
  }

  function cursor_backward() {
    setCurrentPage((currentPage) => Math.max(currentPage - 1, 1));
  }

  function set_data(resp: PaginationResponse) {
    setMaxPage(resp.num_pages);
  }

  return { currentPage, maxPage, cursor_forward, cursor_backward, set_data };
}

export default usePagination;
