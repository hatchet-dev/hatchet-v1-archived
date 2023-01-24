import React from "react";
import { Container, GithubAvatar, GithubName } from "./styles";

type Props = {
  account_name: string;
  account_avatar_url: string;
  avatar_size?: "default" | "small";
};

const GithubAvatarAndName: React.FC<Props> = ({
  account_name,
  account_avatar_url,
  avatar_size = "default",
}) => {
  return (
    <Container height="30px">
      <GithubAvatar src={account_avatar_url} avatar_size={avatar_size} />
      <GithubName>{account_name}</GithubName>
    </Container>
  );
};

export default GithubAvatarAndName;
