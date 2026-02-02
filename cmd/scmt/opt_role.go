package main

import (
	"encoding/json"
	"fmt"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:   messages.GetUse("role"),
	Short: messages.GetShort("role"),
	Long:  messages.GetLong("role"),
}

var roleAddCmd = &cobra.Command{
	Use:   "add <role>",
	Short: "Add a role to the server",
	Long:  "Add a new role to the server's role list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		role := args[0]
		
		cfg := config.New()

		d, err := data.New(*cfg)
		if err != nil {
			return err
		}

		err = d.Open()
		if err != nil {
			return err
		}

		changed, err := d.AddRole(role, Engineer, Message)
		if err != nil {
			return err
		}

		if !changed {
			if OutputJSON {
				output := map[string]interface{}{
					"action":  "add",
					"role":    role,
					"changed": false,
					"message": "Role already exists",
				}
				jsonBytes, _ := json.MarshalIndent(output, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("Role '%s' already exists\n", role)
			}
			return nil
		}

		err = d.Save()
		if err != nil {
			return err
		}

		if OutputJSON {
			output := map[string]interface{}{
				"action":  "add",
				"role":    role,
				"changed": true,
				"message": "Role added successfully",
			}
			jsonBytes, _ := json.MarshalIndent(output, "", "  ")
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Printf("Role '%s' added successfully\n", role)
		}

		return nil
	},
}

var roleRemoveCmd = &cobra.Command{
	Use:   "remove <role>",
	Short: "Remove a role from the server",
	Long:  "Remove a role from the server's role list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		role := args[0]
		
		cfg := config.New()

		d, err := data.New(*cfg)
		if err != nil {
			return err
		}

		err = d.Open()
		if err != nil {
			return err
		}

		changed, err := d.RemoveRole(role, Engineer, Message)
		if err != nil {
			if OutputJSON {
				output := map[string]interface{}{
					"action": "remove",
					"role":   role,
					"error":  err.Error(),
				}
				jsonBytes, _ := json.MarshalIndent(output, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("Error: %s\n", err.Error())
			}
			return nil
		}

		if changed {
			err = d.Save()
			if err != nil {
				return err
			}
		}

		if OutputJSON {
			output := map[string]interface{}{
				"action":  "remove",
				"role":    role,
				"changed": changed,
				"message": "Role removed successfully",
			}
			jsonBytes, _ := json.MarshalIndent(output, "", "  ")
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Printf("Role '%s' removed successfully\n", role)
		}

		return nil
	},
}

var roleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all roles on the server",
	Long:  "Display all roles currently assigned to the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.New()

		d, err := data.New(*cfg)
		if err != nil {
			return err
		}

		err = d.Open()
		if err != nil {
			return err
		}

		roles := d.ListRoles()

		if OutputJSON {
			output := map[string]interface{}{
				"action": "list",
				"roles":  roles,
				"count":  len(roles),
			}
			jsonBytes, _ := json.MarshalIndent(output, "", "  ")
			fmt.Println(string(jsonBytes))
		} else {
			if len(roles) == 0 {
				fmt.Println("No roles assigned to this server")
			} else {
				fmt.Printf("Server roles (%d):\n", len(roles))
				for _, role := range roles {
					fmt.Printf("  - %s\n", role)
				}
			}
		}

		return nil
	},
}

func init() {
	log.Debugf("role command init, start")
	defer log.Debugf("role command init, end")

	roleCmd.AddCommand(roleAddCmd)
	roleCmd.AddCommand(roleRemoveCmd)
	roleCmd.AddCommand(roleListCmd)
	
	rootCmd.AddCommand(roleCmd)
}