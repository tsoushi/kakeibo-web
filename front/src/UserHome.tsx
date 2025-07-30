import { useQuery } from "urql";
import { graphql } from "./gql";
import { Link as RouterLink } from "react-router-dom";
import { Box, Typography, Button, CircularProgress, Stack, Paper } from "@mui/material";

const GetUserDocument = graphql(/* GraphQL */ `
  query GetUser {
    user {
      id
      name
    }
  }
`);

export default function UserHome() {
  const [{ data, fetching, error }] = useQuery({ query: GetUserDocument });

  if (fetching) {
    return (
      <Box p={6} textAlign="center">
        <CircularProgress />
        <Typography mt={2}>Loading...</Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Paper elevation={3} sx={{ p: 3, maxWidth: 400, mx: "auto", textAlign: "center" }}>
        <Typography variant="h6" color="error">Error</Typography>
        <Typography mt={2}>{error.message}</Typography>
        <Button component={RouterLink} to="/record" variant="contained" color="primary" sx={{ mt: 2 }}>Go to Record</Button>
      </Paper>
    );
  }

  return (
    <Paper elevation={3} sx={{ p: 3, maxWidth: 400, mx: "auto" }}>
      <Stack spacing={2} alignItems="center">
        <Typography variant="h5">User Home</Typography>
        <Typography>Welcome to {data?.user.name}'s home page!</Typography>
        <Button component={RouterLink} to="/record" variant="contained" color="success">Go to Record</Button>
      </Stack>
    </Paper>
  );
}
