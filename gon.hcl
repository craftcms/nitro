source = ["./nitro"]
bundle_id = "com.craftcms.nitro"

apple_id {
  username = "support@craftcms.com"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Craft CMS"
}

dmg {
  output_path = "nitro.dmg"
  volume_name = "Nitro"
}

zip {
  output_path = "nitro.zip"
}
