import React from "react";
import { ErrorBoundary } from "react-error-boundary";
import {
  FlexColCenter,
  Placeholder,
  SmallSpan,
} from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

type Props = {
  children?: React.ReactNode;
};

const HatchetErrorBoundary: React.FC<Props> = ({ children }) => {
  const handleError = (err: Error) => {
    console.log(err);
  };

  return (
    <ErrorBoundary onError={handleError} FallbackComponent={ErrorPage}>
      {children}
    </ErrorBoundary>
  );
};

export default HatchetErrorBoundary;

const ErrorPage = ({ error }: any) => (
  <Placeholder>
    <FlexColVerticalCenter height="100%" width="100%">
      <SmallSpan>
        An unexpected error occurred. Please refresh the page and try again.
      </SmallSpan>
    </FlexColVerticalCenter>
  </Placeholder>
);

const FlexColVerticalCenter = styled(FlexColCenter)`
  justify-content: center;
  align-items: center;
`;
