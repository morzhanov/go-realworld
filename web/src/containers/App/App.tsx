import "./App.css";
import {QueryClientProvider, QueryClient} from "react-query";

import AppRouter from "../../router/AppRouter";
import {AppProps} from "./App.interface";

const queryClient = new QueryClient({defaultOptions: {queries: {refetchOnWindowFocus: false}}});

export default function App({history}: AppProps) {
  return (
    <div className="App">
        <QueryClientProvider client={queryClient}>
            <AppRouter history={history} />
        </QueryClientProvider>
    </div>
  );
}
