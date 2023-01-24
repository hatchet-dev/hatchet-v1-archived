import {
  FlexCol,
  FlexRowRight,
  Grid,
  H1,
  H2,
  HorizontalSpacer,
  P,
  SmallSpan,
  Span,
  Breadcrumbs,
  GridCard,
  TextInput,
  SectionArea,
  StandardButton,
  ErrorBar,
  CopyCodeline,
  FlexRowLeft,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { CreatePATResponse } from "shared/api/generated/data-contracts";
import styled from "styled-components";

type Props = {
  account_name: string;
  account_avatar_url: string;
};

const GithubAvatarAndName: React.FC<Props> = ({
  account_name,
  account_avatar_url,
}) => {
  return (
    <Container height="30px">
      <GithubAvatar src={account_avatar_url} />
      <GithubName>{account_name}</GithubName>
    </Container>
  );
};

export default GithubAvatarAndName;

const GithubAvatar = styled.img`
  height: 24px;
  width: 24px;
  border-radius: 100%;
`;

const GithubName = styled(SmallSpan)`
  font-weight: bold;
`;

const Container = styled(FlexRowLeft)`
  gap: 8px;
`;
