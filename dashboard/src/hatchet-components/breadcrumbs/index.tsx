import React from "react";
import { Breadcrumb, BreadcrumbArrow, BreadcrumbWrapper } from "./styles";

export type Breadcrumb = {
  label: string;
  link: string;
};

export type Props = {
  breadcrumbs: Breadcrumb[];
};

const Breadcrumbs: React.FC<Props> = ({ breadcrumbs }) => {
  return (
    <BreadcrumbWrapper>
      {breadcrumbs.map((val, i) => {
        if (i != breadcrumbs.length - 1) {
          return (
            <>
              <Breadcrumb clickable={true} href={val.link}>
                {val.label}
              </Breadcrumb>
              <BreadcrumbArrow className="material-icons">
                keyboard_arrow_right
              </BreadcrumbArrow>
            </>
          );
        }

        return <Breadcrumb clickable={false}>{val.label}</Breadcrumb>;
      })}
    </BreadcrumbWrapper>
  );
};

export default Breadcrumbs;
