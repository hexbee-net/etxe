import { Box } from './Box'

export const Body = ({ overlayColor, children, ...rest }) => (
  <Box
    {...rest}
    flex="1 1 auto"
    display="flex"
    alignItems="center"
    justifyContent="center"
    flexDirection="column"
    bg={overlayColor}
    pe="auto"
  >
    {children}
  </Box>
);
