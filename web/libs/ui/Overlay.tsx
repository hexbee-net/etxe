import styled from "styled-components";
import { animated } from "react-spring";

export const Overlay = styled(animated.div)`
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  pointer-events: none;
  z-index: ${0};
`;
