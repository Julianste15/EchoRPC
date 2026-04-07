package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	pb "servidor.local/grpc-servidor/serviciosCancion"
)

// gRPC Server address
const serverAddr = "localhost:50053"

func main() {
	// 1. Establish connection to the server using gRPC Dial
	// Using WithInsecure() for development connections without TLS
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Could not connect: %v\n", err)
		return
	}
	defer conn.Close()

	// 2. Create gRPC client for the audio service
	c := pb.NewServiciosCancionesClient(conn)

	// 3. Request audio title from user via console
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nEnter the audio title to search for: ")
	tituloLeido, _ := reader.ReadString('\n')
	tituloLeido = strings.TrimSpace(tituloLeido)

	// 4. Prepare the request DTO object
	objPeticion := &pb.PeticionDTO{Titulo: tituloLeido}

	// 5. Configure a context with a 5-second timeout for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 6. Perform the remote call to BuscarAudio procedure
	res, err := c.BuscarAudio(ctx, objPeticion)
	if err != nil {
		// Specific gRPC error handling
		if st, ok := grpcStatus.FromError(err); ok && st.Code() == codes.NotFound {
			fmt.Printf("\nCode: 404 Message: %s\n", st.Message())
			return
		}
		fmt.Printf("Error in gRPC call: %v", err)
		return
	}

	// 7. Process and display the successful response from the server
	fmt.Printf("\nMessage: %s", res.Mensaje)
	fmt.Printf("\nCode: %d", res.Codigo)
	if res.Codigo == 200 {
		fmt.Printf("\n--- Audio Data ---")
		fmt.Printf("\nTitle: %s\nDuration: %d seconds\nType: %s\nAvailable: %v\n",
			res.ObjAudio.Titulo,
			res.ObjAudio.Duracion,
			res.ObjAudio.Tipo,
			res.ObjAudio.Disponible,
		)
	}
	fmt.Println()
}
