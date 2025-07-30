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
        <Button component={RouterLink} to="/" variant="contained" color="primary" sx={{ mt: 2 }}>ホームに戻る</Button>
      </Paper>
    );
  }

  // 現在の年月を取得
  const now = new Date();
  const currentYear = now.getFullYear();
  const currentMonth = now.getMonth() + 1;

  return (
    <Paper elevation={3} sx={{ p: 4, maxWidth: 600, mx: "auto" }}>
      <Stack spacing={3} alignItems="center">
        <Typography variant="h4">家計簿アプリ</Typography>
        <Typography variant="h6">Welcome, {data?.user.name}!</Typography>
        
        <Box sx={{ width: '100%' }}>
          <Typography variant="subtitle1" gutterBottom fontWeight="bold">レコード</Typography>
          <Button 
            component={RouterLink} 
            to={`/record/monthly/${currentYear}/${currentMonth}`} 
            variant="contained" 
            color="primary" 
            fullWidth
          >
            今月のレコード
          </Button>
        </Box>        <Box sx={{ width: '100%' }}>
          <Typography variant="subtitle1" gutterBottom fontWeight="bold">資産管理</Typography>
          <Stack direction="row" spacing={2}>
            <Button 
              component={RouterLink} 
              to="/asset" 
              variant="contained" 
              color="success"
              fullWidth
            >
              資産一覧
            </Button>
            <Button 
              component={RouterLink} 
              to="/asset/category" 
              variant="outlined" 
              color="success"
              fullWidth
            >
              資産カテゴリ
            </Button>
          </Stack>
        </Box>

        <Box sx={{ width: '100%' }}>
          <Typography variant="subtitle1" gutterBottom fontWeight="bold">タグ管理</Typography>
          <Button 
            component={RouterLink} 
            to="/tag" 
            variant="contained" 
            color="secondary"
            fullWidth
          >
            タグ一覧
          </Button>
        </Box>
      </Stack>
    </Paper>
  );
}
